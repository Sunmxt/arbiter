# arbiter

*Derived from project [**CrossMesh**](https://github.com/Sunmxt/utt)*


[![golang](https://img.shields.io/badge/golang->%3D1.9-blue)](https://golang.org/)
[![Build Status](https://travis-ci.com/Sunmxt/arbiter.svg?branch=master)](https://travis-ci.com/Sunmxt/arbiter)
[![codecov](https://codecov.io/gh/Sunmxt/arbiter/branch/master/graph/badge.svg)](https://codecov.io/gh/Sunmxt/arbiter)

## Intro

Arbiter is tracer to **manage goroutine lifecycles** , **preventing risk of goroutine leak**. It also simplifies implement of **Graceful Termination**.



## Getting Started

```go
import arbit "github.com/sunmxt/utt/arbiter"
```

```go
arbiter := arbit.New() // new arbiter
```



### Spawn a goroutine (like "go" keyword)

```go
arbiter.Go(func(){
  // ... do something ...
})
```

### Trace an execution

```go
arbiter.Do(func(){
  // ... do something ...
})
```



### Shutdown

```go
arbiter.Shutdown() // shutdown. Arbiter will send exit signal to all goroutines and executions. 
```

### Intergrate Shutdown() with OS signals

```go
arbiter.StopOSSignals(syscall.SIGTERM, syscall.SIGINT) // SIGTERM and SIGINT will tigger Shutdown().
```

### Watch a shutdown

*Shutdown* signal will be sent via a channel.

```go
select {
  case <-arbiter.Exit(): // watch for a shutdown signal.
    // ...do cleanning...
  case ...
  case ...
}
```

Or you may periodically check **arbiter.Shutdown()**. For example: 

```go
for arbiter.ShouldRun() {
  // ... do something ...
}
// ...do cleanning...
```



### Join (Wait)

```go
arbiter.Join() // blocked until all goroutines and executions exited.
```

#### Arbit

```go
arbiter.Arbit() // Let SIGTERM and SIGINT tigger Shutdown() than wait.
```

*This is an shortcut of:*

```go
arbiter.StopOSSignals(syscall.SIGTERM, syscall.SIGINT)
arbiter.Join()
```



---

### Arbiter tree

Create derived Arbiter.

```go
child := arbit.NewWithParent(arbiter) // new child arbiter
```

Many derived arbiters forms a **arbiter tree**, which has following properties:

- Derived arbiters will be **automatically shut down** when the parent does.
- *Arbiter.Join()* waits for all goroutines and executions on the arbiter tree (i.e **childrens' included** ) to exit



