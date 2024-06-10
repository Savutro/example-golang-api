# Golang-API

## Installationen

Um dieses Repo nutzen zu können muss zwingend Golang version 1.21 installiert sein.

Ausserdem muss MySQL installiert sein. (Mit MySQL Workbench CE kann Einsicht auf die Datenbank genommen werden)

Diese App nutzt Zweifaktorenauthentifizierung im Loginprozess. Dafür muss der GoogleAuthenticator installiert werden.

## VSCode Konfiguration

Für Codereferenzen kann die Erweiterung "Golang" installiert werden.

## Nutzen der API

Zunächst sollten der Befehl `$ go get -v -u` ausgeführt werden um die Dependencies zu aktualisieren. Danach mittels `$ go mod tidy` wird das *go.mod File* aktualisiert und die Checksum neu erstellt. 

Mit dem Befehl `$ go run main.go` wird die Main-Funktion gestartet. In der Konsole sollte durch Logs ersichtlich sein, das der Service läuft.

Die API sollte nun bereit sein, so dass unser Angular Frontend mit diesem kommunizieren kann.

Damit das Frontend jedoch funktioniert muss das Zertifikat **server.crt** im Browser installiert werden damit die TLS Verbindung möglich ist.

Zu dem muss noch das separat mitgelieferte **config.yml** in den config Ordner gepackt werden da so die Umgebungsvariablen mittels **Viper** geladen werden.

> Aus der eben erwähnten Config sind auch die Credentials für den MySql Account zu abzuleiten.

Die App kann auch durch ein Programm wie **Insomnia** oder **Postman** getestet werden. 


## Authentifikations-Flow

Der Flow sollte auch im requests.http File ersichtlich sein

1. `/register` -> User wird in DB erstellt -> Gibt QR Code zurück (qr.png im Repo wenn mit Insomnia getestet wird)
2. `/login` -> Logindaten von Registration übernehmen und einloggen -> Secure Session Cookie wird erstellt. 
3. `/twofactor` -> Code von der Google Authenticator App mitgeben -> Validiert Token in Session Cookie
4. Nun können alle "Book" Endpunkte angesprochen werden. 
5. Der Endpunkt `/admin` ist nur ansprechbar wenn man dem User manuell die Rolle des "admin" erteilt (in der DB)
6. `/logout` -> Der Token im Session Cookie wird invalid gemacht -> Neues Login möglich

## Injection

In **auth.go** kann man den Code auf den **Zeilen 44 - 55** einkommentieren und dafür den Code darunter zwischen den **Zeilen 57 - 64** auskommentieren. Somit ersetzt man für dieses Beispiel die Version mit dem ORM GORM mit dessen bei der ein roher SQL Query gemacht wird.

Wenn man nun beim Registrieren folgenden Input macht (anstelle eines normalen Username):
```
(example', 'password', 'secret', 'admin'); -- 
```
hat man nun einen admin user erstellt. Dies kann man in der DB nachweisen.

## Wahlthemen

Wie mit Herrn Dumermuth besprochen haben wir mehrere Themen behandelt darunter:

- Zweifaktoren Authentifizierung via **GoogleAuth**
- Verschlüsselung der Datenübertragung mittels **Self Signed Cert** (TLS)
- Sanitization via **GORM**
- Passwortmanagement mit **Hashing** und **Salting**
- Autorisierung der Endpunkte mittels **Middlewares**
- Eigenes Softwaredesign basierend auf dem **MVC Modell** 

## Design

Die Struktur basiert auf einer modifizierten Variante eines **MVC Modells**

Im **main.go** File sind die Server Initialisierung und die Routen definiert.

Im **package auth** sind die Helperfunktionen für die Authentifizierung definiert.

Im **package config** befindet sich die definition der Datenbank und die config in Form eines YAML.

Im **package middlewares** sind alle wiederverwendbaren Handler definiert, welche vor den Controllern zwischengeschaltet sind.

Im **package controllers** sind die Handler definiert, welche die Routen bedienen. 

Im **package models** sind die Objekte definiert und deren Methoden oder Helperfunktionen die sich auf diese Beziehen.

Im **package utils** befinden sich alle restlichen Helperfunktionen.

## Sicherheitsrelevante Aspekte/Vektoren

Es ist klar, dass nicht alles abgedeckt wurde. Somit wurden keine Validierungen der Werte vorgenommen, um ein Injectionbeispiel zu haben. Ausserdem fehlen genauere Logs welche klare Meldungen von sich geben.

## Entwicklungsprozess

Die API wurde teilweise schon durch Vorarbeit sehr früh fertiggestellt. Hierbei machte man sich Gedanken bezüglich des Techstacks der verfügbar ist und somit wurde die Entscheidung getroffen, dass Fabian den Fokus auf die API setzt. Benjamin nahm dabei eine unterstützende Rolle bezüglich des Loginflows ein und half beim Testing des Backends. Dadurch das die Arbeit somit limitiert ausfiel konnte er parallel beim Vorbereiten des Pitch vollen Einsatz leisten.

Im Entwicklungsprozess wurden viel Eigenerfahrung aus dem Betrieb als auch sehr gelerntes eingesetzt. 

Die Schnittstelle wurde im Anschluss im Zusammenspiel mit der Entwicklung der UI nochmals optimiert und Fehler konnten so gefunden und behoben werden.

In der letzten Testphase funktionierte das Produkt.

## Nutzen / Relevanz

Die App dient zur Demonstration des Gelernten bezüglich Sicherheitsthematiken. Die Bücherapi hat nicht mehr als den Zweck zur Schaustellung der Autorisierung für diverse Endpunkte.