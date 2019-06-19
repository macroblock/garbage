#ifndef __QMK_H__
#define __QMK_H__

#include <stdint.h>

#include "keyboard.h"

////////////////////////////
// tmk_core/common/matrix.h
////////////////////////////

#if (MATRIX_COLS <= 8)
typedef  uint8_t    matrix_row_t;
#elif (MATRIX_COLS <= 16)
typedef  uint16_t   matrix_row_t;
#elif (MATRIX_COLS <= 32)
typedef  uint32_t   matrix_row_t;
#else
#error "MATRIX_COLS: invalid value"
#endif

#if (MATRIX_ROWS <= 8)
typedef  uint8_t    matrix_col_t;
#elif (MATRIX_ROWS <= 16)
typedef  uint16_t   matrix_col_t;
#elif (MATRIX_ROWS <= 32)
typedef  uint32_t   matrix_col_t;
#else
#error "MATRIX_ROWS: invalid value"
#endif

extern matrix_row_t matrix[MATRIX_ROWS]; //debounced values

#define MATRIX_IS_ON(row, col)  (matrix[row] && (1<<col))
#define MATRIX_SET(row, col) (matrix[row] |= (1<<col))
#define MATRIX_CLEAR(row, col) (matrix[row] &= ~(1<<col))

//////////////////////////////
// tmk_core/common/keyboard.h
//////////////////////////////

/* key matrix position */
typedef struct {
    uint8_t col;
    uint8_t row;
} keypos_t;

/* key event */
typedef struct {
    keypos_t key;
    bool     pressed;
    uint16_t time;
} keyevent_t;

////////////////////////////
// tmk_core/common/action.h
////////////////////////////

/* tapping count and state */
typedef struct {
    bool    interrupted :1;
    bool    reserved2   :1;
    bool    reserved1   :1;
    bool    reserved0   :1;
    uint8_t count       :4;
} tap_t;

/* Key event container for recording */
typedef struct {
    keyevent_t  event;
#ifndef NO_ACTION_TAPPING
    tap_t tap;
#endif
} keyrecord_t;


////////////////////
// additional staff
////////////////////
extern keypos_t scancode_table[256];

void init_scancode_table(int cols, int rows, ...);


void qmk_init_keyboard(void);

#endif // __QMK_H__