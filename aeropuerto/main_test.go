package main

import (
    "fmt"
    "testing"
    "time"
)

type TestResult struct {
    Name     string
    Duration time.Duration
    Config   Config
}

func TestConfiguraciones(t *testing.T) {
    results := make([]TestResult, 0)
    
    fmt.Println("\n=== INICIANDO PRUEBAS DE DIFERENTES CONFIGURACIONES ===")

    baseConfig := Config{
        NumAviones:        10,
        NumPistas:         2,
        NumPuertas:        3,
        MaxColaAviones:    5,
        TiempoBaseControl: 2 * time.Second,
        TiempoBasePista:   3 * time.Second,
        TiempoPuerta:      4 * time.Second,
        VariacionTiempo:   0.25,
    }

    // Test 1: Configuración Base
    fmt.Println("\n=== TEST 1: CONFIGURACIÓN BASE ===")
    printConfig("Configuración Base", baseConfig)
    duration := simularAeropuerto(baseConfig)
    results = append(results, TestResult{"Configuración Base", duration, baseConfig})
    fmt.Printf("Duración: %v\n", duration)

    // Test 2: Cola de aviones duplicada
    fmt.Println("\n=== TEST 2: COLA DUPLICADA ===")
    configColaDuplicada := baseConfig
    configColaDuplicada.MaxColaAviones *= 2
    printConfig("Cola Duplicada", configColaDuplicada)
    duration = simularAeropuerto(configColaDuplicada)
    results = append(results, TestResult{"Cola Duplicada", duration, configColaDuplicada})
    fmt.Printf("Duración: %v\n", duration)

    // Test 3: Variación de tiempo aumentada
    fmt.Println("\n=== TEST 3: VARIACIÓN TIEMPO 25% ===")
    configVariacion := baseConfig
    configVariacion.VariacionTiempo = 0.25
    printConfig("Variación Tiempo 25%", configVariacion)
    duration = simularAeropuerto(configVariacion)
    results = append(results, TestResult{"Variación Tiempo 25%", duration, configVariacion})
    fmt.Printf("Duración: %v\n", duration)

    // Test 4: Cola doble y variación 25%
    fmt.Println("\n=== TEST 4: COLA DOBLE + VARIACIÓN 25% ===")
    configMixta := baseConfig
    configMixta.MaxColaAviones *= 2
    configMixta.VariacionTiempo = 0.25
    printConfig("Cola Doble + Variación 25%", configMixta)
    duration = simularAeropuerto(configMixta)
    results = append(results, TestResult{"Cola Doble + Variación 25%", duration, configMixta})
    fmt.Printf("Duración: %v\n", duration)

    // Test 5: Pistas multiplicadas por 5
    fmt.Println("\n=== TEST 5: PISTAS x5 ===")
    configPistas := baseConfig
    configPistas.NumPistas *= 5
    printConfig("Pistas x5", configPistas)
    duration = simularAeropuerto(configPistas)
    results = append(results, TestResult{"Pistas x5", duration, configPistas})
    fmt.Printf("\n- DURACIÓN: %v\n", duration)

    // Test 6: Pistas x5 con tiempo x5
    fmt.Println("\n=== TEST 6: PISTAS x5 + TIEMPO x5 ===")
    configPistasTiempo := baseConfig
    configPistasTiempo.NumPistas *= 5
    configPistasTiempo.TiempoBasePista *= 5
    printConfig("Pistas x5 + Tiempo x5", configPistasTiempo)
    duration = simularAeropuerto(configPistasTiempo)
    results = append(results, TestResult{"Pistas x5 + Tiempo x5", duration, configPistasTiempo})
    fmt.Printf("Duración: %v\n", duration)

    fmt.Println("\n=== RESUMEN DE PRUEBAS ===")
    for _, result := range results {
        fmt.Printf("%s: %v\n", result.Name, result.Duration)
    }
}

func printConfig(name string, config Config) {
    fmt.Printf("Configuración %s:\n"+
        "- Aviones: %d\n"+
        "- Pistas: %d\n"+
        "- Puertas: %d\n"+
        "- Max Cola: %d\n"+
        "- Variación Tiempo: %.2f%%\n",
        name,
        config.NumAviones,
        config.NumPistas,
        config.NumPuertas,
        config.MaxColaAviones,
        config.VariacionTiempo*100)
}