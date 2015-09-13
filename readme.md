# Content Proxy Server

This is an experimental content proxy server designed specifically to
automatically optimize any graphic content for mobile devices.

The code has two versions developed in Python and Go respectively where either
are considered as a prototype and has limited features.

## Python Version

The Python version is made for Python 2.7 and Pillow with **libjpeg**. The goal
of this version is aimed to optimize, resize and/or crop intelligently JPEG images.

The code is located at the root of this repository.

### Setup

```bash
pip install -r requirements.txt
```

In the newer version of PIP, this will install dependencies just for the executing user.

### How to Run the Service

```bash
python server.py
```

By default, it is listening on port **9500**. Run with `-h` for more information.

## Go Version

The Go version is made for Go 1.4 (and maybe 1.5) and custom libraries and
frameworks. The goal of this version is to do at least what the Python version
can do and emphasize on speed and content optimization.

The code is located at `/go`.

### Setup

```bash
make
```

This `make` will install dependencies to `/go/lib`.

### How to Run the Service

```bash
./server
```

By default, it is listening on port **9500**.
