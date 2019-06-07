
#include <cstdio>
#include <conio.h>
#include <stdint.h>

using namespace std;

#define ID_KEY (0)
#define ID_SIZE (1)
#define ID_MACRO (2)
#define ID_JUMP (3)

#define BITMASK (uint16_t(0x03fff))

#define SET_SIZE(val) ((val&BITMASK)^(ID_SIZE<<14))
#define SET_MACRO(val) ((val&BITMASK)^(ID_MACRO<<14))
#define SET_JUMP(val) ((val&BITMASK)^(ID_JUMP<<14))

#define GET_ID(val) (val>>14)
#define GET_SIZE(val) (val&BITMASK)
#define GET_MACRO(val) (val&BITMASK)
#define GET_JUMP(val) (val&BITMASK)

#define IS_END(val) (val == 0)
#define IS_KEY(val) (val>>14 == ID_KEY)
#define IS_SIZE(val) (val>>14 == SIZE_ID)
#define IS_MACRO(val) (val>>14 == ID_MACRO)
#define IS_JUMP(val) (val>>14 == ID_JUMP)

const uint16_t seq_table[] = { SET_SIZE(3), 'a', '1', 'n', '2', 'z', '3' };
const uint16_t* seq_ptr = &seq_table[0];
uint16_t seq_macro = uint16_t(~0);

#define READ_CODE() (*(seq_ptr++))
#define RESET_CODE() (seq_ptr = &seq_table[0])

bool seq_process(uint16_t keycode) {
    uint16_t const* table = &seq_table[0];

    RESET_CODE();
    uint16_t val = READ_CODE();

    switch (GET_ID(val)) {
    case ID_MACRO:
        seq_macro = GET_MACRO(val);
        break;
    case ID_JUMP:
        printf("### error: JUMP found\n");
        break;
    case ID_SIZE:
    default:
        uint16_t l = 0;
        uint16_t r = GET_SIZE(val)-1;
        printf("debug size -> %d\n", r);
        while (l <= r) {
            uint16_t m = (l+r);
            val = seq_ptr[m<<1];
            printf("debug -> %c\n", val);
            if (val < keycode) {
                l = m+1;
            } else if (val > keycode) {
                r = m-1;
            } else {
                seq_ptr = &seq_ptr[(m<<1)+1];
                val = READ_CODE();
                printf("-> %c\n", val);
                RESET_CODE();
                return true;
            }
        }
        return false;
    }
    return false;
}

bool process_record_user(uint16_t keycode) {
    if (seq_process(keycode)) {
        return false;
    }
    return true;
}

int main()
{
    printf("type key\n");
    while (1) {
        uint16_t key = getch();
        if (key == 0x1b) {
            break;
        }
        if (process_record_user(key)) {
            printf("unhandled keycode: %c\n", key);
        }
    }
}