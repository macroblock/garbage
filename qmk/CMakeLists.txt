cmake_minimum_required(VERSION 3.10)

project(qmk-test)

# set(CMAKE_CXX_COMPILER "g++")
# set(CMAKE_BINARY_DIR $CMAKE_CXX_FLAGS)

set(PROJECT_SOURCE_DIR ${CMAKE_SOURCE_DIR}/sources)
# set(PROJECT_BINARY_DIR ${CMAKE_SOURCE_DIR}/bin)

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Werror")

find_package(Curses REQUIRED)
include_directories(${CURSES_INCLUDE_DIR})

set(SRC_LIST ${PROJECT_SOURCE_DIR}/main.cpp)

add_executable(${PROJECT_NAME} ${SRC_LIST})
target_link_libraries(${PROJECT_NAME} ${CURSES_LIBRARIES})

# get_target_property(EXE_NAME ${PROJECT_NAME} OUTPUT_NAME)
# message("#########" ${EXE_NAME})

add_custom_target(run
    COMMAND bin
    DEPENDS bin
    WORKING_DIRECTORY ${EXECUTABLE_OUTPUT_PATH}
)