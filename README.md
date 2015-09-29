# flagon
Auto-propulate Go structs into command-line flags.

Allows these:

```go
type ChipotleOrder struct {
    Rice                string
    Beans               string
    FajitaVegetables    boolean
    Meat                string
    Salsa               []string
    Corn                boolean
    SourCream           boolean
    Cheese              boolean
    Guacamole           boolean
}

type PlacedOrder struct {
    ChipotleOrder
    Address
    TimeReady     time.Time
}
```

...to accept these:

```bash
$ ./main -placed-order \
    -rice brown \
    -beans pinto \
    -fajita-vegetables \
    -meat barbacoa \
    -salsa mild,hot \
    -guacamole \
    -time-ready 7:56PM \
    -address 'last used'
```

...without you having to write all those flags yourself.
