cmake_minimum_required(VERSION 3.10)

project(shumai LANGUAGES CXX C)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(flashlight CONFIG REQUIRED)

# Flashlight <> NAPI JS bindings
# Enable strict handling of warnings
add_definitions("-Wall -Wextra")
add_definitions(-DNAPI_VERSION=6)
include_directories(${CMAKE_JS_INC})
file(GLOB SOURCE_FILES
    "test/out/*.h"
    "test/out/*.cc")
add_library(flashlight_napi_bindings SHARED ${SOURCE_FILES})
set_target_properties(flashlight_napi_bindings PROPERTIES
    PREFIX ""
    SUFFIX ".node")
target_link_libraries(flashlight_napi_bindings PRIVATE ${CMAKE_JS_LIB})
target_include_directories(flashlight_napi_bindings PRIVATE
    ${CMAKE_SOURCE_DIR}/node_modules/node-addon-api
    ${CMAKE_SOURCE_DIR}/node_modules/node-addon-api/src
    ${CMAKE_JS_INC})
execute_process(COMMAND node -p "require('node-addon-api').include"
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
        OUTPUT_VARIABLE NODE_ADDON_API_DIR
        )
string(REPLACE "\n" "" NODE_ADDON_API_DIR ${NODE_ADDON_API_DIR})
string(REPLACE "\"" "" NODE_ADDON_API_DIR ${NODE_ADDON_API_DIR})
target_include_directories(flashlight_napi_bindings PRIVATE ${NODE_ADDON_API_DIR})

target_link_libraries(
  flashlight_napi_bindings
  PRIVATE 
  flashlight::flashlight
)

