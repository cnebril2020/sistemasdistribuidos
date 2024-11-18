package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// Configuración del sistema
type Config struct {
    NumAviones        int
    NumPistas         int
    NumPuertas        int
    MaxColaAviones    int
    TiempoBaseControl time.Duration
    TiempoBasePista   time.Duration
    TiempoPuerta      time.Duration
    VariacionTiempo   float64
}

// Estructuras principales
type TorreControl struct {
    cola     chan int
    maxCola  int
    ocupada  bool
    mu       sync.Mutex
}

type Pista struct {
    id       int
    ocupada  bool
    mu       sync.Mutex
}

type PuertaEmbarque struct {
    id       int
    ocupada  bool
    mu       sync.Mutex
}

// Función principal de simulación
func simularAeropuerto(config Config) time.Duration {
    startTime := time.Now()
    
    // Inicializar canales y estructuras
    torreControl := &TorreControl{
        cola:    make(chan int, config.MaxColaAviones),
        maxCola: config.MaxColaAviones,
    }

    pistas := make([]*Pista, config.NumPistas)
    for i := range pistas {
        pistas[i] = &Pista{id: i}
    }

    puertas := make([]*PuertaEmbarque, config.NumPuertas)
    for i := range puertas {
        puertas[i] = &PuertaEmbarque{id: i}
    }

    var wg sync.WaitGroup

    // Iniciar aviones
    for i := 0; i < config.NumAviones; i++ {
        wg.Add(1)
        go func(idAvion int) {
            defer wg.Done()
            procesarAvion(idAvion, torreControl, pistas, puertas, config)
        }(i)
    }

    wg.Wait()
    return time.Since(startTime)
}

func procesarAvion(id int, torre *TorreControl, pistas []*Pista, puertas []*PuertaEmbarque, config Config) {
    fmt.Printf("Avión %d solicitando permiso a la torre\n", id)
    
    if err := solicitarTorre(id, torre, config); err != nil {
        fmt.Printf("Avión %d no pudo ser atendido por la torre: %v\n", id, err)
        return
    }

    fmt.Printf("Avión %d ha sido autorizado por la torre\n", id)

    pista := buscarPistaDisponible(pistas)
    if pista == nil {
        fmt.Printf("Avión %d no encontró pista disponible\n", id)
        return
    }

    aterrizar(id, pista, config)

    puerta := buscarPuertaDisponible(puertas)
    if puerta == nil {
        fmt.Printf("Avión %d no encontró puerta disponible\n", id)
        return
    }

    desembarcar(id, puerta, config)
}

func solicitarTorre(id int, torre *TorreControl, config Config) error {
    select {
    case torre.cola <- id:
        torre.mu.Lock()
        torre.ocupada = true
        torre.mu.Unlock()
        
        tiempo := aplicarVariacion(config.TiempoBaseControl, config.VariacionTiempo)
        time.Sleep(tiempo)
        
        <-torre.cola
        torre.mu.Lock()
        torre.ocupada = false
        torre.mu.Unlock()
        return nil
    default:
        return fmt.Errorf("cola de torre llena")
    }
}

func buscarPistaDisponible(pistas []*Pista) *Pista {
    for {
        for _, pista := range pistas {
            pista.mu.Lock()
            if !pista.ocupada {
                pista.ocupada = true
                pista.mu.Unlock()
                return pista
            }
            pista.mu.Unlock()
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func aterrizar(id int, pista *Pista, config Config) {
    fmt.Printf("Avión %d comenzando aterrizaje en pista %d\n", id, pista.id)
    tiempo := aplicarVariacion(config.TiempoBasePista, config.VariacionTiempo)
    time.Sleep(tiempo)
    
    pista.mu.Lock()
    pista.ocupada = false
    pista.mu.Unlock()
    
    fmt.Printf("Avión %d ha aterrizado en pista %d\n", id, pista.id)
}

func buscarPuertaDisponible(puertas []*PuertaEmbarque) *PuertaEmbarque {
    for {
        for _, puerta := range puertas {
            puerta.mu.Lock()
            if !puerta.ocupada {
                puerta.ocupada = true
                puerta.mu.Unlock()
                return puerta
            }
            puerta.mu.Unlock()
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func desembarcar(id int, puerta *PuertaEmbarque, config Config) {
    fmt.Printf("Avión %d iniciando desembarque en puerta %d\n", id, puerta.id)
    tiempo := aplicarVariacion(config.TiempoPuerta, config.VariacionTiempo)
    time.Sleep(tiempo)
    
    puerta.mu.Lock()
    puerta.ocupada = false
    puerta.mu.Unlock()
    
    fmt.Printf("Desembarco del avión %d completado en puerta %d\n", id, puerta.id)
}

func aplicarVariacion(tiempo time.Duration, variacion float64) time.Duration {
    factor := 1.0 + (rand.Float64()*variacion*2 - variacion)
    return time.Duration(float64(tiempo) * factor)
}

func main() {
    rand.Seed(time.Now().UnixNano())

    config := Config{
        NumAviones:        10,
        NumPistas:         2,
        NumPuertas:        3,
        MaxColaAviones:    5,
        TiempoBaseControl: 2 * time.Second,
        TiempoBasePista:   3 * time.Second,
        TiempoPuerta:      4 * time.Second,
        VariacionTiempo:   0.25,
    }

        fmt.Println("Iniciando simulación normal del aeropuerto...")
        fmt.Printf("Configuración:\n- Aviones: %d\n- Pistas: %d\n- Puertas: %d\n- Max Cola: %d\n",
            config.NumAviones, config.NumPistas, config.NumPuertas, config.MaxColaAviones)
        
        duracion := simularAeropuerto(config)
        
        fmt.Printf("Simulación completada en %v\n", duracion)
} 