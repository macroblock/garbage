# cmake version to be used
cmake_minimum_required( VERSION 3.10 )
project( qmk VERSION 0.0.1 LANGUAGES CXX )

find_package( SDL2 REQUIRED )

add_executable( qmk-test ./src/main.cpp ./src/bitmapfont.cpp )
target_include_directories( qmk-test PRIVATE ./include ./src )
target_compile_options( qmk-test PRIVATE -Wall -Wextra -pedantic -Werror )

# external libs
target_link_libraries( qmk-test PRIVATE SDL2 SDL2_image )