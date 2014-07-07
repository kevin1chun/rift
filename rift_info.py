from flask import Flask
from flask import json
import service

app = Flask(__name__)

@app.route('/about')
def about():
    return json.dumps({
        '@context': 'http://schema.org',
        '@type': 'SoftwareApplication',
        'name': 'Rift',
        'image': 'http://cdn29.elitedaily.com/wp-content/uploads/2013/06/main-iStock_000000292864Small1.jpg',
        'url': 'https://github.com/jerluc/rift',
        'description': 'A protocol for decentralized service distribution'
    })

desc = {
    '@context': 'http://schema.org',
    '@type': 'DiscoverAction',
    'name': 'Rift',
    'description': 'What is Rift?',
    'target': {
        '@type': 'EntryPoint',
        'urlTemplate': 'http://172.16.1.52:8677/about',
        'contentType': 'application/json+ld',
        'httpMethod': 'GET'
    }
}

if __name__ == '__main__':
    service.register(desc)
    app.run(host='0.0.0.0', port=8677)
