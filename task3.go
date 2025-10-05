package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	A          = 0.0
	B          = math.Pi / 2.0
	Accuracy   = 0.5e-4
	ScanSteps  = 2000
	WorkerPool = 0
	MaxIter    = 100000
	KRelax     = 0.1
)

type Result struct {
	intervalA float64
	intervalB float64
	mid       float64

	bisectionRoot      float64
	bisectionIters     int
	bisectionFval      float64
	bisectionConverged bool

	iterRoot      float64
	iterIters     int
	iterFval      float64
	iterConverged bool

	errorIterVsBisection float64
}

func f(x float64) float64 {
	return x*x*math.Cos(2*x) + 1.0
}

func scanner(ctx context.Context, out chan<- [2]float64, steps int) {
	defer close(out)
	step := (B - A) / float64(steps)
	a := A
	for i := 0; i < steps; i++ {
		b := a + step
		fa := f(a)
		fb := f(b)
		if math.IsNaN(fa) || math.IsNaN(fb) {
			a = b
			continue
		}
		if fa*fb <= 0 {
			select {
			case out <- [2]float64{a, b}:
			case <-ctx.Done():
				return
			}
		}
		a = b
	}
}

func bisection(a, b float64, tol float64, maxIter int) (root float64, iters int, fval float64, converged bool) {
	fa := f(a)
	fb := f(b)
	if fa*fb > 0 {
		return math.NaN(), 0, math.NaN(), false
	}
	left, right := a, b
	var mid float64
	for iters = 1; iters <= maxIter; iters++ {
		mid = 0.5 * (left + right)
		fm := f(mid)
		if math.IsNaN(fm) {
			return mid, iters, fm, false
		}
		if math.Abs(fm) <= tol || (right-left)/2 <= tol {
			return mid, iters, fm, true
		}
		if fa*fm <= 0 {
			right = mid
			fb = fm
		} else {
			left = mid
			fa = fm
		}
	}
	return mid, iters, f(mid), false
}

func iteration(x0 float64, k float64, tol float64, maxIter int) (root float64, iters int, fval float64, converged bool) {
	x := x0
	for iters = 1; iters <= maxIter; iters++ {
		fx := f(x)
		xn := x - k*fx
		if math.IsNaN(fx) || math.IsNaN(xn) || math.IsInf(xn, 0) {
			return xn, iters, fx, false
		}
		if math.Abs(xn-x) <= tol || math.Abs(fx) <= tol {
			return xn, iters, f(xn), true
		}
		x = xn
	}
	return x, iters, f(x), false
}

func worker(ctx context.Context, id int, in <-chan [2]float64, out chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("worker %d got panic: %v\n", id, r)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case interval, ok := <-in:
			if !ok {
				return
			}
			a := interval[0]
			b := interval[1]
			mid := 0.5 * (a + b)

			res := Result{intervalA: a, intervalB: b, mid: mid}

			rootB, itB, fB, convB := bisection(a, b, Accuracy, MaxIter)
			res.bisectionRoot = rootB
			res.bisectionIters = itB
			res.bisectionFval = fB
			res.bisectionConverged = convB

			rootI, itI, fI, convI := iteration(mid, KRelax, Accuracy, MaxIter)
			res.iterRoot = rootI
			res.iterIters = itI
			res.iterFval = fI
			res.iterConverged = convI

			if !math.IsNaN(rootB) && !math.IsNaN(rootI) {
				res.errorIterVsBisection = math.Abs(rootI - rootB)
			} else {
				res.errorIterVsBisection = math.NaN()
			}

			select {
			case out <- res:
			case <-ctx.Done():
				return
			}
		}
	}
}

func saveResultsCSV(filename string, results []Result) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	head := []string{"a", "b", "mid", "method", "root", "f(root)", "iters", "converged", "error_vs_bisection"}
	if err = w.Write(head); err != nil {
		return err
	}

	for _, r := range results {
		rowB := []string{
			fmt.Sprintf("%g", r.intervalA),
			fmt.Sprintf("%g", r.intervalB),
			fmt.Sprintf("%g", r.mid),
			"bisection",
			fmt.Sprintf("%.12g", r.bisectionRoot),
			fmt.Sprintf("%.12g", r.bisectionFval),
			fmt.Sprintf("%d", r.bisectionIters),
			fmt.Sprintf("%t", r.bisectionConverged),
			"",
		}
		if err = w.Write(rowB); err != nil {
			return err
		}

		rowI := []string{
			fmt.Sprintf("%g", r.intervalA),
			fmt.Sprintf("%g", r.intervalB),
			fmt.Sprintf("%g", r.mid),
			"iteration",
			fmt.Sprintf("%.12g", r.iterRoot),
			fmt.Sprintf("%.12g", r.iterFval),
			fmt.Sprintf("%d", r.iterIters),
			fmt.Sprintf("%t", r.iterConverged),
			fmt.Sprintf("%.12g", r.errorIterVsBisection),
		}
		if err = w.Write(rowI); err != nil {
			return err
		}
	}
	return nil
}

/*
Створити багатопоточний (з розпаралеленими обчисленнями) проект (використовуючи горутини та канали)
для наближеного розв'язання нелінійного рівняння (відповідно до варіанту (див. таблицю 4)) заданими методами,
метод сканування використовувати для знаходження відрізків локалізації,
інші методи використовувати для ітераційного уточнення коріння.
Виконати порівняльний аналіз за точністю обчислень (якщо можливо) та за кількістю ітерацій.
Результат обчислень вивести у файл. Використовувати помилки, паніки та відновлення
*/
func main() {
	start := time.Now()
	if WorkerPool <= 0 {
		WorkerPool = runtime.NumCPU()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	intervals := make(chan [2]float64, 128)
	resultsCh := make(chan Result, 128)

	go scanner(ctx, intervals, ScanSteps)

	var wg sync.WaitGroup
	for i := 0; i < WorkerPool; i++ {
		wg.Add(1)
		go worker(ctx, i, intervals, resultsCh, &wg)
	}

	collected := make([]Result, 0)
	var collectWg sync.WaitGroup
	collectWg.Add(1)
	go func() {
		defer collectWg.Done()
		for r := range resultsCh {
			collected = append(collected, r)
		}
	}()

	wg.Wait()
	close(resultsCh)
	collectWg.Wait()

	outFilename := "results.csv"
	if err := saveResultsCSV(outFilename, collected); err != nil {
		fmt.Printf("Failed to save results: %v\n", err)
		return
	}

	var bIterSum, iIterSum, bCount, iCount int
	var bConvCount, iConvCount int
	for _, r := range collected {
		if r.bisectionIters > 0 {
			bIterSum += r.bisectionIters
			bCount++
		}
		if r.iterIters > 0 {
			iIterSum += r.iterIters
			iCount++
		}
		if r.bisectionConverged {
			bConvCount++
		}
		if r.iterConverged {
			iConvCount++
		}
	}

	fmt.Println("Done calculating")
	fmt.Printf("Found %d intervals (scanned %d subintervals).\n", len(collected), ScanSteps)
	if bCount > 0 {
		fmt.Printf("Bisection: avg iters = %.2f, converged %d/%d\n", float64(bIterSum)/float64(bCount), bConvCount, bCount)
	}
	if iCount > 0 {
		fmt.Printf("Iteration: avg iters = %.2f, converged %d/%d\n", float64(iIterSum)/float64(iCount), iConvCount, iCount)
	}

	for idx, r := range collected {
		fmt.Printf("root %d: interval=[%.6g, %.6g], b_root=%.10g (iters=%d, conv=%t), i_root=%.10g (iters=%d, conv=%t), err=%.12g\n",
			idx+1, r.intervalA, r.intervalB, r.bisectionRoot, r.bisectionIters, r.bisectionConverged, r.iterRoot, r.iterIters, r.iterConverged, r.errorIterVsBisection)
	}

	duration := time.Since(start)
	fmt.Printf("Elapsed: %s\n", duration)
}
