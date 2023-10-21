package SAKnife

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Middleware para verificar los roles
func VerificarRolesMiddleware(rolesRequeridos []string, siguiente http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener el token JWT del encabezado de autorización
		bearer_string := "Bearer"
		tokenString := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), bearer_string, "", -1))

		// Verificar y decodificar el token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		fmt.Println(token)
		// Verificar los roles del usuario a partir del token JWT
		claims, ok := token.Claims.(jwt.MapClaims)
		fmt.Println(ok)

		if !ok {
			http.Error(w, "Error al obtener los roles del token", http.StatusUnauthorized)
			return
		}
		fmt.Println(claims)
		sub, ok := claims["sub"]
		fmt.Println(sub)

		roles, ok := claims["role"]
		aroles := roles.(string)
		aroles = fmt.Sprintf(",%s,", aroles)
		troles := strings.Split(aroles, ",")
		fmt.Println(troles)
		if !ok {
			http.Error(w, "Error al obtener los roles del token", http.StatusUnauthorized)
			return
		}

		// Verificar si el usuario tiene el rol necesario
		haRequeridoRol := false
		for _, rol := range rolesRequeridos {
			if strings.Contains(aroles, rol) {
				haRequeridoRol = true
				break
			}

		}

		if !haRequeridoRol {
			http.Error(w, "Acceso no autorizado", http.StatusForbidden)
			return
		}

		// Si la verificación de roles pasa, continúa con el siguiente manejador
		siguiente(w, r)
	}
}

// Función para requerir un rol específico
func RequiereRol(roles string) []string {
	return strings.Split(roles, ",")

}
