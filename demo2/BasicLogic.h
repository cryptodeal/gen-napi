#include <stdint.h>
#include <string>

namespace demo2 {

int8_t foo(int8_t a);

double* bar(double* a, int32_t b);

float* baz(float* a, int b);

long long* qux(long long* a, int b);

bool* quux(bool a, bool b);

std::string test(std::string a);

double* test2(std::vector<double> a);

void* test3(std::vector<void*> a);
}  // namespace demo2