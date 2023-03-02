#include <stdint.h>
#include <string>

namespace demo2 {

void foo(int8_t a);

std::vector<double> bar(double* a, int32_t b);

std::vector<float> baz(float* a, int b);

double qux(long long* a, int b);

std::vector<bool> quux(bool a, bool b);

std::string test(std::string a);

std::array<int8_t, 20> test2(std::vector<double> a);

std::pair<int, int> test3(double* a);

char* test4(const char* a, const char16_t* b);

std::vector<long long> test5(std::vector<double> a);

}  // namespace demo2