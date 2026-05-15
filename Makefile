BIN=prefwatch
CODESIGN_IDENTITY=Developer ID Application: JESSE GORDON DONAT (NBWN497MH2)
NOTARY_PROFILE=notarytool-profile

.PHONY: all
all: clean build

.PHONY: test
test:
	go test ./...

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	-rm -rf release dist
	mkdir release dist

release/darwin_amd64/$(BIN):
	mkdir -p release/darwin_amd64
	env GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o release/darwin_amd64/$(BIN)

release/darwin_arm64/$(BIN):
	mkdir -p release/darwin_arm64
	env GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o release/darwin_arm64/$(BIN)

release/darwin_universal/$(BIN): release/darwin_amd64/$(BIN) release/darwin_arm64/$(BIN)
	mkdir -p release/darwin_universal
	lipo -create -output release/darwin_universal/$(BIN) release/darwin_amd64/$(BIN) release/darwin_arm64/$(BIN)

.PHONY: build
build: release/darwin_universal/$(BIN)

.PHONY: sign
sign: build
	codesign \
		--force \
		--timestamp \
		--options runtime \
		--sign "$(CODESIGN_IDENTITY)" \
		release/darwin_universal/$(BIN)

	codesign --verify --strict --verbose=4 release/darwin_universal/$(BIN)

.PHONY: package
package: sign
	mkdir -p dist
	ditto -c -k --keepParent release/darwin_universal/$(BIN) dist/$(BIN).darwin_universal.zip

.PHONY: notarize
notarize: package
	xcrun notarytool submit dist/$(BIN).darwin_universal.zip \
		--keychain-profile "$(NOTARY_PROFILE)" \
		--wait

.PHONY: release
release: clean notarize
