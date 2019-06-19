
#include <stdint.h>
#include <cstdarg>

#include "qmk.h"


matrix_row_t matrix[MATRIX_ROWS]; //debounced values

keypos_t scancode_table[256];

void init_scancode_table(int cols, int rows, ...)
{
    for(int i = 0; i <= 255; i++) {
        scancode_table[i] = keypos_t{255, 255};
    }
    auto n_args = cols*rows;

    va_list keys;
    va_start(keys, rows);
    for(int i = 0; i < n_args; i++) {
        auto c = i % cols;
        auto r = i / cols;
        if (r > rows-1) {
            break;
        }
        int k = va_arg(keys, int);
        if (k > 255) { continue; }
        scancode_table[k] = keypos_t{uint8_t(c), uint8_t(r)};
    }
    va_end(keys);
}