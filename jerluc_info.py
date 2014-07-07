from flask import Flask
from flask import json
import service

app = Flask(__name__)

@app.route('/about')
def about():
    return json.dumps({
        '@context': 'http://schema.org',
        '@type': 'Person',
        'name': 'Jeremy Lucas',
        'image': 'http://urx.com/assets/team/jeremy-38a148756b5598cbd08fc67a26ff9900.jpg',
        'url': 'http://jerluc.com',
        'description': 'A really weird guy',
    })

desc = {
    '@context': 'http://schema.org',
    '@type': 'DiscoverAction',
    'name': 'Jeremy Lucas',
    'description': 'Who the hell is he?',
    'target': {
        '@type': 'EntryPoint',
        'urlTemplate': 'http://172.16.1.52:8676/about',
        'contentType': 'application/json+ld',
        'httpMethod': 'GET'
    }
}

if __name__ == '__main__':
    service.register(desc)
    app.run(host='0.0.0.0', port=8676)
