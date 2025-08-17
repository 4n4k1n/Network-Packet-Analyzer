# Geographic IP Location Implementation Guide

This guide walks you through adding geographic IP location tracking to your Network Packet Analyzer, allowing you to see where network traffic originates from around the world.

## Table of Contents
1. [Overview](#overview)
2. [Setup and Installation](#setup-and-installation)
3. [Implementation Steps](#implementation-steps)
4. [Code Examples](#code-examples)
5. [Performance Optimization](#performance-optimization)
6. [Output Examples](#output-examples)
7. [Troubleshooting](#troubleshooting)
8. [Advanced Features](#advanced-features)

## Overview

Geographic IP (GeoIP) lookup allows you to determine the physical location of an IP address, including:
- Country and city
- Latitude and longitude coordinates
- Internet Service Provider (ISP)
- Autonomous System Number (ASN)
- Organization information

### Why Add GeoIP?
- **Security Analysis**: Identify suspicious traffic from unexpected countries
- **Performance Monitoring**: Track latency to different geographic regions
- **Compliance**: Monitor data flows for regulatory requirements
- **Network Optimization**: Understand traffic patterns by location

## Setup and Installation

### Step 1: Install Required Dependencies

```bash
# Install GeoIP2 Go library
go get github.com/oschwald/geoip2-golang

# Update your go.mod
go mod tidy
```

### Step 2: Obtain GeoIP Database

#### Option A: MaxMind GeoLite2 (Recommended - Free)

1. **Register for free account**: https://www.maxmind.com/en/geolite2/signup
2. **Generate license key**: Account → My License Key → Generate new license key
3. **Download database**:
   ```bash
   # Create directory for GeoIP data
   mkdir -p ./data/geoip
   
   # Download GeoLite2-City database (replace YOUR_LICENSE_KEY)
   wget "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=YOUR_LICENSE_KEY&suffix=tar.gz" -O GeoLite2-City.tar.gz
   
   # Extract database
   tar -xzf GeoLite2-City.tar.gz
   mv GeoLite2-City_*/GeoLite2-City.mmdb ./data/geoip/
   ```

#### Option B: Alternative Free APIs

```bash
# No setup required, but rate limited:
# - ip-api.com: 1000 requests/hour
# - ipinfo.io: 1000 requests/day
# - freegeoip.app: 15000 requests/hour
```

## Implementation Steps

### Step 1: Update Data Structures

Add to `src/types.go`:

```go
import (
    "sync"
    "time"
)

// Geographic location information
type GeoLocation struct {
    Country     string  `json:"country"`
    CountryCode string  `json:"country_code"`
    City        string  `json:"city"`
    Latitude    float64 `json:"latitude"`
    Longitude   float64 `json:"longitude"`
    ISP         string  `json:"isp"`
    ASN         string  `json:"asn"`
    Timezone    string  `json:"timezone"`
}

// Enhanced packet data with geographic info
type Data struct {
    src_ip       string
    dst_ip       string
    protocol     layers.IPProtocol
    tcp          *layers.TCP
    udp          *layers.UDP
    src_port     uint16
    dst_port     uint16
    src_geo      GeoLocation  // NEW: Source IP location
    dst_geo      GeoLocation  // NEW: Destination IP location
    packet_size  int          // NEW: Packet size in bytes
    timestamp    time.Time    // NEW: Packet timestamp
}

// GeoIP cache for performance
type GeoCache struct {
    cache map[string]GeoLocation
    mutex sync.RWMutex
    hits  int
    misses int
}

// Enhanced statistics with geographic data
type Stats_data struct {
    // ... existing fields ...
    geo_cache        *GeoCache
    countries        map[string]int  // Country packet counts
    cities          map[string]int  // City packet counts
    international_packets int        // Non-local traffic count
    local_packets    int            // Local network traffic count
}
```

### Step 2: Create GeoIP Module

Create `src/geoip.go`:

```go
package main

import (
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "strings"
    "sync"
    "time"
    
    "github.com/oschwald/geoip2-golang"
)

// Global GeoIP database and cache
var (
    geoipDB    *geoip2.Reader
    geoCache   *GeoCache
    useLocalDB bool = true  // Switch between local DB and API
)

// Initialize GeoIP functionality
func initGeoIP() error {
    // Initialize cache
    geoCache = &GeoCache{
        cache: make(map[string]GeoLocation),
        mutex: sync.RWMutex{},
    }
    
    // Try to open local database
    var err error
    geoipDB, err = geoip2.Open("./data/geoip/GeoLite2-City.mmdb")
    if err != nil {
        fmt.Printf("Warning: Could not open GeoIP database: %v\n", err)
        fmt.Println("Falling back to online API (rate limited)")
        useLocalDB = false
        return nil
    }
    
    fmt.Println("GeoIP database loaded successfully")
    useLocalDB = true
    return nil
}

// Cleanup GeoIP resources
func closeGeoIP() {
    if geoipDB != nil {
        geoipDB.Close()
    }
}

// Check if IP is private/local
func isPrivateIP(ip string) bool {
    parsedIP := net.ParseIP(ip)
    if parsedIP == nil {
        return false
    }
    
    // Check for private IP ranges
    private := []string{
        "10.0.0.0/8",
        "172.16.0.0/12", 
        "192.168.0.0/16",
        "127.0.0.0/8",
        "169.254.0.0/16",
        "::1/128",
        "fc00::/7",
        "fe80::/10",
    }
    
    for _, cidr := range private {
        _, subnet, _ := net.ParseCIDR(cidr)
        if subnet.Contains(parsedIP) {
            return true
        }
    }
    return false
}

// Get geographic location for IP address
func getGeoLocation(ip string) GeoLocation {
    // Return empty location for private IPs
    if isPrivateIP(ip) {
        return GeoLocation{
            Country: "Local",
            City:    "Private Network",
        }
    }
    
    // Check cache first
    geoCache.mutex.RLock()
    if cached, exists := geoCache.cache[ip]; exists {
        geoCache.mutex.RUnlock()
        geoCache.hits++
        return cached
    }
    geoCache.mutex.RUnlock()
    geoCache.misses++
    
    var geo GeoLocation
    
    if useLocalDB {
        geo = getGeoFromLocalDB(ip)
    } else {
        geo = getGeoFromAPI(ip)
    }
    
    // Cache the result
    geoCache.mutex.Lock()
    geoCache.cache[ip] = geo
    geoCache.mutex.Unlock()
    
    return geo
}

// Lookup using local MaxMind database
func getGeoFromLocalDB(ip string) GeoLocation {
    parsedIP := net.ParseIP(ip)
    if parsedIP == nil {
        return GeoLocation{Country: "Invalid IP"}
    }
    
    record, err := geoipDB.City(parsedIP)
    if err != nil {
        return GeoLocation{Country: "Unknown"}
    }
    
    return GeoLocation{
        Country:     record.Country.Names["en"],
        CountryCode: record.Country.IsoCode,
        City:        record.City.Names["en"],
        Latitude:    record.Location.Latitude,
        Longitude:   record.Location.Longitude,
        Timezone:    record.Location.TimeZone,
        ASN:         fmt.Sprintf("AS%d", record.Traits.AutonomousSystemNumber),
    }
}

// Lookup using online API (rate limited)
func getGeoFromAPI(ip string) GeoLocation {
    // Using ip-api.com (free tier: 1000 requests/hour)
    url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,city,lat,lon,timezone,isp,as", ip)
    
    client := &http.Client{Timeout: 2 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return GeoLocation{Country: "API Error"}
    }
    defer resp.Body.Close()
    
    var apiResp struct {
        Status      string  `json:"status"`
        Country     string  `json:"country"`
        CountryCode string  `json:"countryCode"`
        City        string  `json:"city"`
        Lat         float64 `json:"lat"`
        Lon         float64 `json:"lon"`
        Timezone    string  `json:"timezone"`
        ISP         string  `json:"isp"`
        AS          string  `json:"as"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return GeoLocation{Country: "Parse Error"}
    }
    
    if apiResp.Status != "success" {
        return GeoLocation{Country: "Unknown"}
    }
    
    return GeoLocation{
        Country:     apiResp.Country,
        CountryCode: apiResp.CountryCode,
        City:        apiResp.City,
        Latitude:    apiResp.Lat,
        Longitude:   apiResp.Lon,
        Timezone:    apiResp.Timezone,
        ISP:         apiResp.ISP,
        ASN:         apiResp.AS,
    }
}

// Get cache statistics
func getCacheStats() (int, int, float64) {
    total := geoCache.hits + geoCache.misses
    hitRate := 0.0
    if total > 0 {
        hitRate = float64(geoCache.hits) / float64(total) * 100
    }
    return geoCache.hits, geoCache.misses, hitRate
}
```

### Step 3: Update Main Program

Modify `src/main.go`:

```go
func main() {
    // ... existing initialization ...
    
    // Initialize GeoIP
    if err := initGeoIP(); err != nil {
        log.Printf("GeoIP initialization failed: %v", err)
    }
    defer closeGeoIP()
    
    // ... existing packet capture loop ...
    
    // Enhanced data collection in packet loop:
    for pack := range pack_src.Packets() {
        stats_data.total_packets++
        data = getData(pack)
        
        // Add geographic lookups (async to avoid blocking)
        go func(d *Data) {
            d.src_geo = getGeoLocation(d.src_ip)
            d.dst_geo = getGeoLocation(d.dst_ip)
        }(&data)
        
        printPacketData(data)
        updateGeoStats(&stats_data, data)
        
        // ... rest of existing loop ...
    }
    
    printStats(stats_data)
    printGeoStats(stats_data)
}
```

### Step 4: Update Output Functions

Modify `src/prints.go`:

```go
// Enhanced packet display with geographic info
func printPacketData(data Data) {
    srcGeo := fmt.Sprintf("%s, %s", data.src_geo.City, data.src_geo.Country)
    dstGeo := fmt.Sprintf("%s, %s", data.dst_geo.City, data.dst_geo.Country)
    
    if srcGeo == ", " { srcGeo = "Local" }
    if dstGeo == ", " { dstGeo = "Local" }
    
    fmt.Printf("%-15s %-20s %-15s %-20s %-8s\n", 
        data.src_ip, srcGeo,
        data.dst_ip, dstGeo,
        data.protocol)
}

// New header for geographic display
func printHeaderLine() {
    fmt.Printf("%-15s %-20s %-15s %-20s %-8s\n", 
        "Source IP", "Source Location", 
        "Dest IP", "Dest Location", 
        "Protocol")
}

// Geographic statistics
func printGeoStats(stats_data Stats_data) {
    fmt.Printf("\n=== Geographic Analysis ===\n")
    
    // Top countries
    fmt.Printf("Top Countries:\n")
    topCountries := getTopN(stats_data.countries, 5)
    for i, item := range topCountries {
        percentage := float64(item.count) / float64(stats_data.total_packets) * 100
        fmt.Printf("  %d. %s: %d packets (%.1f%%)\n", 
            i+1, item.key, item.count, percentage)
    }
    
    // Top cities
    fmt.Printf("\nTop Cities:\n")
    topCities := getTopN(stats_data.cities, 5)
    for i, item := range topCities {
        fmt.Printf("  %d. %s: %d packets\n", 
            i+1, item.key, item.count)
    }
    
    // International vs local traffic
    intlPercentage := float64(stats_data.international_packets) / float64(stats_data.total_packets) * 100
    localPercentage := float64(stats_data.local_packets) / float64(stats_data.total_packets) * 100
    
    fmt.Printf("\nTraffic Distribution:\n")
    fmt.Printf("  International: %d packets (%.1f%%)\n", 
        stats_data.international_packets, intlPercentage)
    fmt.Printf("  Local Network: %d packets (%.1f%%)\n", 
        stats_data.local_packets, localPercentage)
    
    // Cache performance
    hits, misses, hitRate := getCacheStats()
    fmt.Printf("\nGeoIP Cache Performance:\n")
    fmt.Printf("  Cache hits: %d\n", hits)
    fmt.Printf("  Cache misses: %d\n", misses)
    fmt.Printf("  Hit rate: %.1f%%\n", hitRate)
}

// Helper function for sorting statistics
type countItem struct {
    key   string
    count int
}

func getTopN(m map[string]int, n int) []countItem {
    items := make([]countItem, 0, len(m))
    for k, v := range m {
        items = append(items, countItem{k, v})
    }
    
    // Simple bubble sort for top N
    for i := 0; i < len(items)-1; i++ {
        for j := 0; j < len(items)-i-1; j++ {
            if items[j].count < items[j+1].count {
                items[j], items[j+1] = items[j+1], items[j]
            }
        }
    }
    
    if n > len(items) {
        n = len(items)
    }
    return items[:n]
}

// Update geographic statistics
func updateGeoStats(stats *Stats_data, data Data) {
    // Initialize maps if needed
    if stats.countries == nil {
        stats.countries = make(map[string]int)
        stats.cities = make(map[string]int)
    }
    
    // Count source locations
    if data.src_geo.Country != "" {
        stats.countries[data.src_geo.Country]++
        if data.src_geo.City != "" {
            cityKey := fmt.Sprintf("%s, %s", data.src_geo.City, data.src_geo.Country)
            stats.cities[cityKey]++
        }
        
        if data.src_geo.Country == "Local" {
            stats.local_packets++
        } else {
            stats.international_packets++
        }
    }
    
    // Count destination locations  
    if data.dst_geo.Country != "" {
        stats.countries[data.dst_geo.Country]++
        if data.dst_geo.City != "" {
            cityKey := fmt.Sprintf("%s, %s", data.dst_geo.City, data.dst_geo.Country)
            stats.cities[cityKey]++
        }
    }
}
```

## Performance Optimization

### 1. Asynchronous Lookups

```go
// Channel for background geo lookups
type GeoLookupJob struct {
    IP     string
    Result chan GeoLocation
}

var geoJobQueue = make(chan GeoLookupJob, 100)

// Start background workers
func startGeoWorkers(numWorkers int) {
    for i := 0; i < numWorkers; i++ {
        go geoWorker()
    }
}

func geoWorker() {
    for job := range geoJobQueue {
        result := getGeoLocation(job.IP)
        select {
        case job.Result <- result:
        case <-time.After(100 * time.Millisecond):
            // Timeout if no one is listening
        }
    }
}
```

### 2. Configurable Lookup Frequency

```go
// Only lookup every Nth unique IP
var lookupCounter = make(map[string]int)

func shouldLookupGeo(ip string) bool {
    lookupCounter[ip]++
    return lookupCounter[ip] == 1 || lookupCounter[ip]%10 == 0
}
```

### 3. Persistent Cache

```go
// Save cache to disk
func saveCacheToFile(filename string) error {
    geoCache.mutex.RLock()
    defer geoCache.mutex.RUnlock()
    
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    return json.NewEncoder(file).Encode(geoCache.cache)
}

// Load cache from disk
func loadCacheFromFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    geoCache.mutex.Lock()
    defer geoCache.mutex.Unlock()
    
    return json.NewDecoder(file).Decode(&geoCache.cache)
}
```

## Output Examples

### Enhanced Real-time Display
```
Source IP       Source Location      Dest IP         Dest Location        Protocol
192.168.1.100   Local               8.8.8.8         Mountain View, US    UDP
192.168.1.100   Local               172.217.16.142  Ashburn, US          TCP
10.0.0.1        Local               151.101.1.140   San Francisco, US    TCP
192.168.1.50    Local               104.244.42.129  Tokyo, Japan         TCP
```

### Geographic Statistics
```
=== Geographic Analysis ===
Top Countries:
  1. United States: 1,247 packets (67.2%)
  2. Germany: 298 packets (16.1%)
  3. Japan: 156 packets (8.4%)
  4. United Kingdom: 89 packets (4.8%)
  5. Canada: 45 packets (2.4%)

Top Cities:
  1. Mountain View, United States: 423 packets
  2. Frankfurt, Germany: 298 packets
  3. Tokyo, Japan: 156 packets
  4. London, United Kingdom: 89 packets
  5. Toronto, Canada: 45 packets

Traffic Distribution:
  International: 608 packets (32.8%)
  Local Network: 1247 packets (67.2%)

GeoIP Cache Performance:
  Cache hits: 1,543
  Cache misses: 234
  Hit rate: 86.8%
```

## Troubleshooting

### Common Issues

**1. "GeoIP database not found"**
```bash
# Verify database path
ls -la ./data/geoip/GeoLite2-City.mmdb

# Re-download if missing
wget "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=YOUR_KEY&suffix=tar.gz"
```

**2. "API rate limit exceeded"**
```bash
# Switch to local database or implement rate limiting
time.Sleep(100 * time.Millisecond) // Between API calls
```

**3. "Slow performance"**
```go
// Enable async lookups and increase cache size
const MAX_CACHE_SIZE = 10000
```

**4. "Memory usage too high"**
```go
// Implement cache eviction
if len(geoCache.cache) > MAX_CACHE_SIZE {
    // Remove oldest entries
    clearOldCacheEntries()
}
```

### Debug Mode

Add debug flags to see lookup details:

```go
var debugGeo = flag.Bool("debug-geo", false, "Enable GeoIP debug output")

func debugPrintf(format string, args ...interface{}) {
    if *debugGeo {
        fmt.Printf("[GEO DEBUG] "+format, args...)
    }
}
```

## Advanced Features

### 1. Export Geographic Data

```go
// Export to CSV with coordinates
func exportGeoCSV(filename string, data []Data) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    writer := csv.NewWriter(file)
    defer writer.Flush()
    
    // Write header
    writer.Write([]string{"src_ip", "src_country", "src_city", "src_lat", "src_lon", 
                          "dst_ip", "dst_country", "dst_city", "dst_lat", "dst_lon"})
    
    // Write data
    for _, d := range data {
        record := []string{
            d.src_ip, d.src_geo.Country, d.src_geo.City, 
            fmt.Sprintf("%.6f", d.src_geo.Latitude), 
            fmt.Sprintf("%.6f", d.src_geo.Longitude),
            d.dst_ip, d.dst_geo.Country, d.dst_geo.City,
            fmt.Sprintf("%.6f", d.dst_geo.Latitude), 
            fmt.Sprintf("%.6f", d.dst_geo.Longitude),
        }
        writer.Write(record)
    }
    
    return nil
}
```

### 2. Geographic Filtering

```go
// Filter by country
-geo-filter "country=US,JP"          // Only US and Japan
-geo-exclude "country=CN,RU"         // Exclude China and Russia
-geo-distance "lat=40.7,lon=-74.0,radius=100"  // Within 100km of NYC
```

### 3. Interactive Map Output

```go
// Generate HTML map with Leaflet.js
func generateInteractiveMap(data []Data) string {
    html := `
    <!DOCTYPE html>
    <html>
    <head>
        <link rel="stylesheet" href="https://unpkg.com/leaflet/dist/leaflet.css" />
        <script src="https://unpkg.com/leaflet/dist/leaflet.js"></script>
    </head>
    <body>
        <div id="map" style="height: 600px;"></div>
        <script>
            var map = L.map('map').setView([40.7, -74.0], 2);
            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png').addTo(map);
    `
    
    // Add markers for each unique location
    for _, d := range data {
        if d.src_geo.Latitude != 0 && d.src_geo.Longitude != 0 {
            html += fmt.Sprintf(`
                L.marker([%.6f, %.6f]).addTo(map)
                    .bindPopup('%s<br>%s');
            `, d.src_geo.Latitude, d.src_geo.Longitude, 
               d.src_geo.City, d.src_ip)
        }
    }
    
    html += `
        </script>
    </body>
    </html>`
    
    return html
}
```

### 4. Security Analysis

```go
// Detect suspicious geographic patterns
func analyzeGeoSecurity(stats Stats_data) {
    suspiciousCountries := []string{"CN", "RU", "KP", "IR"}
    
    fmt.Printf("\n=== Security Analysis ===\n")
    for _, country := range suspiciousCountries {
        if count, exists := stats.countries[country]; exists {
            fmt.Printf("⚠️  Traffic from %s: %d packets\n", country, count)
        }
    }
    
    // Check for rapid country changes
    // Check for impossible geographic hops
    // Alert on first-time countries
}
```

This implementation provides a complete geographic IP tracking system for your packet analyzer, with both offline database support and online API fallback, caching for performance, and rich geographic statistics.

## Next Steps

1. **Start Simple**: Implement basic country lookup first
2. **Add Caching**: Essential for performance with high packet rates  
3. **Enhance Display**: Show geographic info in packet output
4. **Add Statistics**: Country/city breakdown in summary
5. **Optimize**: Async lookups and cache management
6. **Extend**: Interactive maps, security analysis, geographic filtering

Remember to respect API rate limits and consider privacy implications when implementing geographic tracking!