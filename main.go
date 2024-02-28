package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Estructura para mapear la respuesta JSON.
type ApiResponse struct {
	AsnName    string   `json:"asnName"`
	AsnHandle  int64    `json:"asnHandle"`
	OrgID      string   `json:"orgID"`
	OrgName    string   `json:"orgName"`
	OrgCountry string   `json:"orgCountry"`
	Ipv4Prefix []string `json:"ipv4_prefix"`
	Ipv6Prefix []string `json:"ipv6_prefix"`
}

func main() {
	// Definición de flags para los parámetros opcionales
	orgname := flag.String("orgname", "", "Organización (opcional)")
	asn := flag.String("asn", "", "Número de ASN (opcional)")
	ip := flag.String("ip", "", "Dirección IP (opcional)")
	cidr := flag.String("cidr", "", "CIDR (opcional)")
	apiKey := flag.String("apikey", "", "Clave de API para el servicio")

	// Parsear los argumentos
	flag.Parse()

	// Construir la URL de la solicitud
	baseURL := "https://asn-lookup.p.rapidapi.com/api"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Añadir parámetros opcionales si están presentes
	params := url.Values{}
	if *orgname != "" {
		params.Add("orgname", *orgname)
	}
	if *asn != "" {
		params.Add("asn", *asn)
	}
	if *ip != "" {
		params.Add("ip", *ip)
	}
	if *cidr != "" {
		params.Add("cidr", *cidr)
	}
	reqURL.RawQuery = params.Encode()

	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Añadir la cabecera de la clave de API
	req.Header.Add("X-RapidAPI-Key", *apiKey)
	req.Header.Add("X-RapidAPI-Host", "asn-lookup.p.rapidapi.com")

	// Realizar la solicitud
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error performing request:", err)
		return
	}
	defer res.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Deserializar JSON a un array de ApiResponse
	var apiResponses []ApiResponse
	err = json.Unmarshal(body, &apiResponses)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	// Iterar sobre el array e imprimir los contenidos de IPv4Prefix de cada objeto
	for _, apiResp := range apiResponses {
		for _, ipv4Range := range apiResp.Ipv4Prefix {
			fmt.Println(ipv4Range)
		}
	}
}
