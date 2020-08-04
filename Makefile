COMPILER_NAME=ogen_compiler
LAUNCHER_NAME=ogen_launcher
COMPILER_FOLDER=ogen-compiler-linux_amd64
LAUNCHER_FOLDER=ogen-launcher-linux_amd64

build: build_compiler build_launcher

build_compiler:
	cd compiler && go build -o $(COMPILER_NAME) main.go && mv $(COMPILER_NAME) ../

build_launcher:
	cd launcher && go build -o $(LAUNCHER_NAME) main.go && mv $(LAUNCHER_NAME) ../

pack: build pack_compiler pack_launcher

pack_compiler:
	rm -rf $(COMPILER_FOLDER)
	mkdir $(COMPILER_FOLDER) && mv ./$(COMPILER_NAME) $(COMPILER_FOLDER)
	tar -czvf $(COMPILER_FOLDER).tar.gz $(COMPILER_FOLDER)
	rm -rf $(COMPILER_FOLDER)

pack_launcher:
	rm -rf $(LAUNCHER_FOLDER)
	mkdir $(LAUNCHER_FOLDER) && mv ./$(LAUNCHER_NAME) $(LAUNCHER_FOLDER)
	tar -czvf $(LAUNCHER_FOLDER).tar.gz $(LAUNCHER_FOLDER)
	rm -rf $(LAUNCHER_FOLDER)

clean:
	cd compiler && go clean
	cd launcher && go clean
	rm -rf $(COMPILER_NAME) && rm -rf compiler/$(COMPILER_NAME)
	rm -rf $(LAUNCHER_NAME) && rm -rf launcher/$(LAUNCHER_NAME)
	rm -rf *.tar.gz
	rm -rf $(COMPILER_FOLDER)
	rm -rf $(LAUNCHER_FOLDER)