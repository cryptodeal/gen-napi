#include <stdint.h>
#include <string>

namespace demo2 {

enum class dtype {
  f16 = 0,  // 16-bit float
  f32 = 1,  // 32-bit float
  f64 = 2,  // 64-bit float
  b8 = 3,   // 8-bit boolean
  s16 = 4,  // 16-bit signed integer
  s32 = 5,  // 32-bit signed integer
  s64 = 6,  // 64-bit signed integer
  u8 = 7,   // 8-bit unsigned integer
  u16 = 8,  // 16-bit unsigned integer
  u32 = 9,  // 32-bit unsigned integer
  u64 = 10  // 64-bit unsigned integer
};

float* baz(float* a, int b);

std::string test(std::string a);
}