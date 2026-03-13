TLS_DIR := tls
CERT_FILE := $(TLS_DIR)/cert.pem
KEY_FILE := $(TLS_DIR)/key.pem
DAYS ?= 365
CN ?= localhost

.PHONY: tls clean-tls

tls:
	mkdir -p $(TLS_DIR)
	openssl req -x509 -newkey rsa:2048 -nodes \
		-keyout $(KEY_FILE) \
		-out $(CERT_FILE) \
		-days $(DAYS) \
		-subj "/CN=$(CN)"
	@echo "Generated $(CERT_FILE) and $(KEY_FILE)"

clean-tls:
	rm -f $(CERT_FILE) $(KEY_FILE)
