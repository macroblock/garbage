#ifndef __QMK_H__
#define __QMK_H__

#include <stdint.h>

// tmk_core/common/keyboard.h

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

// tmk_core/common/action.h

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

#endif // __QMK_H__