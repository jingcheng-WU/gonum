#!/usr/bin/env bash

# Copyright ©2015 The Gonum Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

WARNINGF32='//\
// Float32 implementations are autogenerated and not directly tested.\
'
WARNINGC64='//\
// Complex64 implementations are autogenerated and not directly tested.\
'

# Level1 routines.

echo Generating level1float32.go
echo -e '// Code generated by "go generate github.com/jingcheng-WU/gonum/blas/gonum”; DO NOT EDIT.\n' > level1float32.go
cat level1float64.go \
| gofmt -r 'blas.Float64Level1 -> blas.Float32Level1' \
\
| gofmt -r 'float64 -> float32' \
| gofmt -r 'blas.DrotmParams -> blas.SrotmParams' \
\
| gofmt -r 'f64.AxpyInc -> f32.AxpyInc' \
| gofmt -r 'f64.AxpyUnitary -> f32.AxpyUnitary' \
| gofmt -r 'f64.DotUnitary -> f32.DotUnitary' \
| gofmt -r 'f64.L2NormInc -> f32.L2NormInc' \
| gofmt -r 'f64.L2NormUnitary -> f32.L2NormUnitary' \
| gofmt -r 'f64.ScalInc -> f32.ScalInc' \
| gofmt -r 'f64.ScalUnitary -> f32.ScalUnitary' \
\
| sed -e "s_^\(func (Implementation) \)D\(.*\)\$_$WARNINGF32\1S\2_" \
      -e 's_^// D_// S_' \
      -e "s_^\(func (Implementation) \)Id\(.*\)\$_$WARNINGF32\1Is\2_" \
      -e 's_^// Id_// Is_' \
      -e 's_"github.com/jingcheng-WU/gonum/internal/asm/f64"_"github.com/jingcheng-WU/gonum/internal/asm/f32"_' \
      -e 's_"math"_math "github.com/jingcheng-WU/gonum/internal/math32"_' \
>> level1float32.go

echo Generating level1cmplx64.go
echo -e '// Code generated by "go generate github.com/jingcheng-WU/gonum/blas/gonum”; DO NOT EDIT.\n' > level1cmplx64.go
cat level1cmplx128.go \
| gofmt -r 'blas.Complex128Level1 -> blas.Complex64Level1' \
\
| gofmt -r 'float64 -> float32' \
| gofmt -r 'complex128 -> complex64' \
\
| gofmt -r 'c128.AxpyInc -> c64.AxpyInc' \
| gofmt -r 'c128.AxpyUnitary -> c64.AxpyUnitary' \
| gofmt -r 'c128.DotcInc -> c64.DotcInc' \
| gofmt -r 'c128.DotcUnitary -> c64.DotcUnitary' \
| gofmt -r 'c128.DotuInc -> c64.DotuInc' \
| gofmt -r 'c128.DotuUnitary -> c64.DotuUnitary' \
| gofmt -r 'c128.ScalInc -> c64.ScalInc' \
| gofmt -r 'c128.ScalUnitary -> c64.ScalUnitary' \
| gofmt -r 'dcabs1 -> scabs1' \
\
| sed -e "s_^\(func (Implementation) \)Zdot\(.*\)\$_$WARNINGC64\1Cdot\2_" \
      -e 's_^// Zdot_// Cdot_' \
      -e "s_^\(func (Implementation) \)Zdscal\(.*\)\$_$WARNINGC64\1Csscal\2_" \
      -e 's_^// Zdscal_// Csscal_' \
      -e "s_^\(func (Implementation) \)Z\(.*\)\$_$WARNINGC64\1C\2_" \
      -e 's_^// Z_// C_' \
      -e "s_^\(func (Implementation) \)Iz\(.*\)\$_$WARNINGC64\1Ic\2_" \
      -e 's_^// Iz_// Ic_' \
      -e "s_^\(func (Implementation) \)Dz\(.*\)\$_$WARNINGC64\1Sc\2_" \
      -e 's_^// Dz_// Sc_' \
      -e 's_"github.com/jingcheng-WU/gonum/internal/asm/c128"_"github.com/jingcheng-WU/gonum/internal/asm/c64"_' \
      -e 's_"math"_math "github.com/jingcheng-WU/gonum/internal/math32"_' \
>> level1cmplx64.go

echo Generating level1float32_sdot.go
echo -e '// Code generated by "go generate github.com/jingcheng-WU/gonum/blas/gonum”; DO NOT EDIT.\n' > level1float32_sdot.go
cat level1float64_ddot.go \
| gofmt -r 'float64 -> float32' \
\
| gofmt -r 'f64.DotInc -> f32.DotInc' \
| gofmt -r 'f64.DotUnitary -> f32.DotUnitary' \
\
| sed -e "s_^\(func (Implementation) \)D\(.*\)\$_$WARNINGF32\1S\2_" \
      -e 's_^// D_// S_' \
      -e 's_"github.com/jingcheng-WU/gonum/internal/asm/f64"_"github.com/jingcheng-WU/gonum/internal/asm/f32"_' \
>> level1float32_sdot.go

echo Generating level1float32_dsdot.go
echo -e '// Code generated by "go generate github.com/jingcheng-WU/gonum/blas/gonum”; DO NOT EDIT.\n' > level1float32_dsdot.go
cat level1float64_ddot.go \
| gofmt -r '[]float64 -> []float32' \
\
| gofmt -r 'f64.DotInc -> f32.DdotInc' \
| gofmt -r 'f64.DotUnitary -> f32.DdotUnitary' \
\
| sed -e "s_^\(func (Implementation) \)D\(.*\)\$_$WARNINGF32\1Ds\2_" \
      -e 's_^// D_// Ds_' \
      -e 's_"github.com/jingcheng-WU/gonum/internal/asm/f64"_"github.com/jingcheng-WU/gonum/internal/asm/f32"_' \
>> level1float32_dsdot.go

echo Generating level1float32_sdsdot.go
echo -e '// Code generated by "go generate github.com/jingcheng-WU/gonum/blas/gonum”; DO NOT EDIT.\n' > level1float32_sdsdot.go
cat level1float64_ddot.go \
| gofmt -r 'float64 -> float32' \
\
| gofmt -r 'f64.DotInc(x, y, f(n), f(incX), f(incY), f(ix), f(iy)) -> alpha + float32(f32.DdotInc(x, y, f(n), f(incX), f(incY), f(ix), f(iy)))' \
| gofmt -r 'f64.DotUnitary(a, b) -> alpha + float32(f32.DdotUnitary(a, b))' \
\
| sed -e "s_^\(func (Implementation) \)D\(.*\)\$_$WARNINGF32\1Sds\2_" \
      -e 's_^// D\(.*\)$_// Sds\1 plus a constant_' \
      -e 's_\\sum_alpha + \\sum_' \
      -e 's/n int/n int, alpha float32/' \
      -e 's_"github.com/jingcheng-WU/gonum/internal/asm/f64"_"github.com/jingcheng-WU/gonum/internal/asm/f32"_' \
>> level1float32_sdsdot.go


# Level2 routines.

echo Generating level2float32.go
echo -e '// Code generated by "go generate github.com/jingcheng-WU/gonum/blas/gonum”; DO NOT EDIT.\n' > level2float32.go
cat level2float64.go \
| gofmt -r 'blas.Float64Level2 -> blas.Float32Level2' \
\
| gofmt -r 'float64 -> float32' \
\
| gofmt -r 'f64.AxpyInc -> f32.AxpyInc' \
| gofmt -r 'f64.AxpyIncTo -> f32.AxpyIncTo' \
| gofmt -r 'f64.AxpyUnitary -> f32.AxpyUnitary' \
| gofmt -r 'f64.AxpyUnitaryTo -> f32.AxpyUnitaryTo' \
| gofmt -r 'f64.DotInc -> f32.DotInc' \
| gofmt -r 'f64.DotUnitary -> f32.DotUnitary' \
| gofmt -r 'f64.ScalInc -> f32.ScalInc' \
| gofmt -r 'f64.ScalUnitary -> f32.ScalUnitary' \
| gofmt -r 'f64.Ger -> f32.Ger' \
\
| sed -e "s_^\(func (Implementation) \)D\(.*\)\$_$WARNINGF32\1S\2_" \
      -e 's_^// D_// S_' \
      -e 's_"github.com/jingcheng-WU/gonum/internal/asm/f64"_"github.com/jingcheng-WU/gonum/internal/asm/f32"_' \
>> level2float32.go

echo Generating level2cmplx64.go
echo -e '// Code generated by "go generate github.com/jingcheng-WU/gonum/blas/gonum”; DO NOT EDIT.\n' > level2cmplx64.go
cat level2cmplx128.go \
| gofmt -r 'blas.Complex128Level2 -> blas.Complex64Level2' \
\
| gofmt -r 'complex128 -> complex64' \
| gofmt -r 'float64 -> float32' \
\
| gofmt -r 'c128.AxpyInc -> c64.AxpyInc' \
| gofmt -r 'c128.AxpyUnitary -> c64.AxpyUnitary' \
| gofmt -r 'c128.DotuInc -> c64.DotuInc' \
| gofmt -r 'c128.DotuUnitary -> c64.DotuUnitary' \
| gofmt -r 'c128.ScalInc -> c64.ScalInc' \
| gofmt -r 'c128.ScalUnitary -> c64.ScalUnitary' \
\
| sed -e "s_^\(func (Implementation) \)Z\(.*\)\$_$WARNINGC64\1C\2_" \
      -e 's_^// Z_// C_' \
      -e 's_"github.com/jingcheng-WU/gonum/internal/asm/c128"_"github.com/jingcheng-WU/gonum/internal/asm/c64"_' \
      -e 's_"math/cmplx"_cmplx "github.com/jingcheng-WU/gonum/internal/cmplx64"_' \
>> level2cmplx64.go

# Level3 routines.

echo Generating level3float32.go
echo -e '// Code generated by "go generate github.com/jingcheng-WU/gonum/blas/gonum”; DO NOT EDIT.\n' > level3float32.go
cat level3float64.go \
| gofmt -r 'blas.Float64Level3 -> blas.Float32Level3' \
\
| gofmt -r 'float64 -> float32' \
\
| gofmt -r 'f64.AxpyUnitaryTo -> f32.AxpyUnitaryTo' \
| gofmt -r 'f64.AxpyUnitary -> f32.AxpyUnitary' \
| gofmt -r 'f64.DotUnitary -> f32.DotUnitary' \
| gofmt -r 'f64.ScalUnitary -> f32.ScalUnitary' \
\
| sed -e "s_^\(func (Implementation) \)D\(.*\)\$_$WARNINGF32\1S\2_" \
      -e 's_^// D_// S_' \
      -e 's_"github.com/jingcheng-WU/gonum/internal/asm/f64"_"github.com/jingcheng-WU/gonum/internal/asm/f32"_' \
>> level3float32.go

echo Generating sgemm.go
echo -e '// Code generated by "go generate github.com/jingcheng-WU/gonum/blas/gonum”; DO NOT EDIT.\n' > sgemm.go
cat dgemm.go \
| gofmt -r 'float64 -> float32' \
| gofmt -r 'sliceView64 -> sliceView32' \
\
| gofmt -r 'dgemmParallel -> sgemmParallel' \
| gofmt -r 'computeNumBlocks64 -> computeNumBlocks32' \
| gofmt -r 'dgemmSerial -> sgemmSerial' \
| gofmt -r 'dgemmSerialNotNot -> sgemmSerialNotNot' \
| gofmt -r 'dgemmSerialTransNot -> sgemmSerialTransNot' \
| gofmt -r 'dgemmSerialNotTrans -> sgemmSerialNotTrans' \
| gofmt -r 'dgemmSerialTransTrans -> sgemmSerialTransTrans' \
\
| gofmt -r 'f64.AxpyInc -> f32.AxpyInc' \
| gofmt -r 'f64.AxpyUnitary -> f32.AxpyUnitary' \
| gofmt -r 'f64.DotUnitary -> f32.DotUnitary' \
\
| sed -e "s_^\(func (Implementation) \)D\(.*\)\$_$WARNINGF32\1S\2_" \
      -e 's_^// D_// S_' \
      -e 's_^// d_// s_' \
      -e 's_"github.com/jingcheng-WU/gonum/internal/asm/f64"_"github.com/jingcheng-WU/gonum/internal/asm/f32"_' \
>> sgemm.go

echo Generating level3cmplx64.go
echo -e '// Code generated by "go generate github.com/jingcheng-WU/gonum/blas/gonum”; DO NOT EDIT.\n' > level3cmplx64.go
cat level3cmplx128.go \
| gofmt -r 'blas.Complex128Level3 -> blas.Complex64Level3' \
\
| gofmt -r 'float64 -> float32' \
| gofmt -r 'complex128 -> complex64' \
\
| gofmt -r 'c128.ScalUnitary -> c64.ScalUnitary' \
| gofmt -r 'c128.DscalUnitary -> c64.SscalUnitary' \
| gofmt -r 'c128.DotcUnitary -> c64.DotcUnitary' \
| gofmt -r 'c128.AxpyUnitary -> c64.AxpyUnitary' \
| gofmt -r 'c128.DotuUnitary -> c64.DotuUnitary' \
\
| sed -e "s_^\(func (Implementation) \)Z\(.*\)\$_$WARNINGC64\1C\2_" \
      -e 's_^// Z_// C_' \
      -e 's_"github.com/jingcheng-WU/gonum/internal/asm/c128"_"github.com/jingcheng-WU/gonum/internal/asm/c64"_' \
      -e 's_"math/cmplx"_cmplx "github.com/jingcheng-WU/gonum/internal/cmplx64"_' \
>> level3cmplx64.go
