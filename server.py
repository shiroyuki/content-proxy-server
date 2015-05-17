from base64   import b64decode
from hashlib  import sha1, md5
import math
import os
import re
from StringIO import StringIO
import sys
from time     import time

from flask import Flask, make_response, render_template
from PIL   import Image
import requests

FLAG_DEBUG_MODE = '--debug'

arguments = sys.argv[1:]

app_name           = 'http-image-optimizer'
image_quality      = 80
disk_cache_enabled = True
debug_mode_enabled = FLAG_DEBUG_MODE in sys.argv
resource_path      = None

remote_cache_map      = {}
remote_cache_lifespan = 3600 # in seconds

app = Flask(__name__)

def main():
    arguments.remove(FLAG_DEBUG_MODE)

    if arguments:
        resource_path = arguments[0]

    #ip   = os.environ['OPENSHIFT_PYTHON_IP']       if 'OPENSHIFT_PYTHON_IP'   in os.environ else '0.0.0.0'
    port = int(os.environ['OPENSHIFT_PYTHON_PORT'] if 'OPENSHIFT_PYTHON_PORT' in os.environ else 9500)

    options = {
        'host':     '0.0.0.0',
        'debug':    debug_mode_enabled,
        'port':     port,
        'threaded': True,
    }

    app.run(**options)

def _reference_path(main_node, *nodes):
    """ File Path """
    if not resource_path:
        raise RuntimeError('The reference resource path is not defined.')

    if main_node[0] == '/':
        return os.path.join(main_node, *nodes)

    return os.path.abspath(os.path.join(resource_path, main_node, *nodes))

def _local_path(main_node, *nodes):
    """ File Path """
    if main_node[0] == '/':
        return os.path.join(main_node, *nodes)

    base_path = os.path.dirname(__file__)

    return os.path.abspath(os.path.join(base_path, main_node, *nodes))

def _hash(content):
    """ Hashing """
    m = sha1()
    m.update(content)

    n = md5()
    n.update(content)

    return '{}x{}'.format(m.hexdigest(), n.hexdigest())

def _process_local_image(path, expected_width, expected_height):
    cache_key  = _hash(path)
    write_path = _get_cache_path(cache_key)

    if _has_cache(write_path):
        return write_path

    return _process_image(path, cache_key, expected_width, expected_height)

def _process_remote_image(url, expected_width, expected_height):
    global remote_cache_map
    global remote_cache_lifespan

    cache_key  = _hash(url)
    write_path = _get_cache_path(cache_key)

    if _has_cache(write_path):
        return write_path

    current_time = time()
    remote_cache = remote_cache_map[url] if url in remote_cache_map else None

    if not remote_cache or remote_cache['expired_at'] < current_time:
        response = requests.get(url)

        if response.status_code != requests.codes.ok:
            raise RuntimeError('Failed to retrieve the image')

        remote_cache_map[url] = {
            'content':    response.content,
            'expired_at': current_time + remote_cache_lifespan
        }

    image_buffer = StringIO(remote_cache_map[url]['content'])

    return _process_image(image_buffer, cache_key, expected_width, expected_height)

def _has_cache(cache_path):
    global disk_cache_enabled

    return disk_cache_enabled and os.path.exists(cache_path)

def _get_cache_path(cache_key):
    cache_path = 'static/cached-images'
    write_path = _local_path('{}/{}.jpg'.format(cache_path, cache_key))

    if not os.path.exists(_local_path(cache_path)):
        os.makedirs(_local_path(cache_path))

    return write_path

def _process_image(source, cache_key, expected_width, expected_height):
    """ Image Processing

        * Not supported variable cropping.
    """
    global app
    global image_quality
    global disk_cache_enabled

    write_path = _get_cache_path(cache_key)

    im = Image.open(source)

    width, height         = im.size
    expected_aspect_ratio = None

    if expected_width and expected_height:
        expected_aspect_ratio = float(expected_width) / float(expected_height)

    original_aspect_ratio = float(width) / float(height)

    if expected_aspect_ratio and original_aspect_ratio != expected_aspect_ratio:
        # Choose the minimum width and height for the cropped region.

        # Calculate for the cropped size.
        w1, h1 = _dimension(width, None, expected_aspect_ratio)
        w2, h2 = _dimension(None, height, expected_aspect_ratio)
        wm, hm = min(w1, w2), min(h1, h2)

        # Calculate for the padding.
        padding_x = (width  - wm) / 2
        padding_y = (height - hm) / 2

        cropped_box = [padding_x, padding_y, width - padding_x, height - padding_y];

        app.logger.debug('Crop: {}'.format(cropped_box))

        im = im.crop(cropped_box)

    # Resize the image
    if expected_width or expected_height:
        if not expected_width:
            expected_width = _dimension(None, expected_height, original_aspect_ratio)[0]

        if not expected_height:
            expected_height = _dimension(expected_width, None, original_aspect_ratio)[1]

        app.logger.debug('Resize: {} x {}'.format(expected_width, expected_height))

        im = im.resize((expected_width, expected_height))

    # Save the colour version.
    im.save(write_path, 'JPEG', quality=image_quality)

    return write_path

def _dimension(width, height, aspect_ratio):
    new_width  = None
    new_height = None

    if not width and height:
        new_width  = int(math.ceil(height * aspect_ratio))

    if not height and width:
        new_height = int(math.ceil(width / aspect_ratio))

    return new_width or width, new_height or height

@app.route("/<width>/<height>/<source_hash>")
def index(width, height, source_hash):
    original_source = b64decode(source_hash)
    cache_key       = _hash('{},{}x{}'.format(original_source, width, height))

    width  = None if width.lower()  == 'auto' else int(width)
    height = None if height.lower() == 'auto' else int(height)

    if re.search('^https?://', original_source):
        # Deal with remote images.
        #raise NotImplemented('Not yet implemented')
        cache_path = _process_remote_image(original_source, width, height)
    else:
        # Deal with local images.
        reference_path = _reference_path(original_source[1:])
        cache_path     = _process_local_image(reference_path, width, height)

    content = None

    with open(cache_path, 'rb') as f:
        content = f.read()

    resp = make_response(content, 200)
    resp.headers['Content-Type'] = 'image/jpeg'

    return resp

if __name__ == "__main__":
    main()