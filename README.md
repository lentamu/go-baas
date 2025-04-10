# go-baas

A lightweight password hashing and verification service written in Go.

## API

**Hash:**

```bash
curl -X POST "http://localhost:8080/hash" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "raw=<password>&cost=<cost>"
```

**Verify:**

```bash
curl -X POST "http://localhost:8080/verify" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "raw=<password>&hash=<hash>"
```
