# TMI Tracks Api

A Golang based backend for the TMI Tracks Suite.
***

### Dev Environment
You can run ```gin --appPort 8080 --port 3000``` to launch a hot reloaded dev environment.

## Paths
***
### ```/create_report```
- POST
- Request format(json):
```
{
  date: string,
  technician: string, // Form YYYY-MM-DD
  notes: string,
  duration: string // Form: HH:mm,
  equipment: string
 }
```
- Response:
```
{
  _id: ID //ID of created object.
}
```
