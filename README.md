# tq
Simple timer queue and a map

_tq_ is basically a map with expiring items, with method called on item expiration. Since all items have the 
relative expiration times, they are processed strictly in order of insertion. Items can be accessed or removed 
at any time, or item expiration time can be refreshed.

## Example
```go
import "github.com/pantonov/tq" 

// This function provides item expiration duration. Note that if value returned
// by this function changes, it does not immediately change order of pending items in timer queue. 
// We use this function (rather than a constant duraiton) to allow on-the-fly configuration changes. 
func dur_func() time.Duration {
    return time.Duration(3 * time.Second) // items will expire in 3 seconds
}

func on_exire(k int, v string) {
	fmt.Printf("Expired item with key %d, value %s at %s\n", 
        k, v, time.Now().Format(time.RFC1123))
}

func main() {
    tq := tq.NewTimerQueue[int, string](on_expire, dur_func) // key type: int, value tyoe: string
    tq.Push(1, "aaa")
    tq.Push(2, "bbb")
    v := tq.Get(1)
    fmt.Printf("item value: %s", *v)
    tq.Remove(1)
    tq.Refresh(2)
    // ... only item '2' will be eventually processed
}	
```

# License
The Unlicense