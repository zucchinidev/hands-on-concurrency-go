## Barycenter problem

<p>
    <img src="https://raw.githubusercontent.com/zucchinidev/hands-on-concurrency-go/master/naiveBarycenter/img/barycenter_problem.png" />
</p>

### Implement a non-concurrent solution to the data parallel barycenter problem.

* Create a program to generate random bodies (cmd folder): Eg.: go run naiveBaryCenter/cmd/main.go 1000000 > naiveBaryCenter/1millbodies.txt
* Load those bodies from a file into memory
* Find the barycenter of those loaded bodies

#### Result 

<p>
    <img src="https://raw.githubusercontent.com/zucchinidev/hands-on-concurrency-go/master/naiveBarycenter/img/result.png" />
</p>

#### Result after refactor

<p>
    <img src="https://raw.githubusercontent.com/zucchinidev/hands-on-concurrency-go/master/naiveBarycenter/img/result_2.png" />
</p>