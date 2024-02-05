import http from 'k6/http'

export let options = {
    vus: 100,
    duration: '5m',
}

export default function(){
    http.get('http://host.docker.internal:8888/new-match')
}