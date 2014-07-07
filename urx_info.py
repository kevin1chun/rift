from flask import Flask
from flask import json
import service

app = Flask(__name__)

@app.route('/about')
def about():
    return json.dumps({
        '@context': 'http://schema.org',
        '@type': 'Organization',
        'name': 'URX',
        'image': 'http://urx.com/assets/random-images/5-db43a60ca3060a09bc2c19e343d1d833.jpg',
        'url': 'http://www.urx.com',
        'description': 'URX is San Francisco-based company that uses mobile deep linking technology to link content across devices. URX works with developers and marketers to reconnect their apps to the web and intelligently route users across devices.',
        'logo': 'http://urx.com/assets/logo-6d896b3ae5dda8477d0dee50fd923805.svg',
        'address': {
            '@type': 'PostalAddress',
            'streetAddress': '168 South Park St.',
            'postalCode': '94107',
            'addressLocality': 'San Francisco',
            'addressRegion': 'California',
            'addressCountry': 'US'
        }
    })

desc = {
    '@context': 'http://schema.org',
    '@type': 'DiscoverAction',
    'name': 'URX',
    'description': 'Learn more about URX',
    'target': {
        '@type': 'EntryPoint',
        'urlTemplate': 'http://172.16.1.52:8675/about',
        'contentType': 'application/json+ld',
        'httpMethod': 'GET'
    }
}

if __name__ == '__main__':
    service.register(desc)
    app.run(host='0.0.0.0', port=8675)
