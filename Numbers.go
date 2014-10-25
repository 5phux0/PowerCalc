package pcalc

import (
	"fmt"
	"math"
)

var primes []int64

type SdigitsFloat struct {
	n  float64
	sd uint8 //max 15
}

type Fraction struct {
	n, d, om int64
}

type number interface {
	FractionValue() *Fraction //returns nil when called on SdigitsFloat
	FloatValue() *SdigitsFloat
	Description() string
}

//Number interface functions
func (f *Fraction) Description() string {
	if f.d != 1 && f.om != 1 {
		return fmt.Sprintf("%d/%d = %d/%d = %f", f.n*f.om, f.d*f.om, f.n, f.d, float64(f.n)/float64(f.d))
	} else if f.d != 1 {
		return fmt.Sprintf("%d/%d = %f", f.n, f.d, float64(f.n)/float64(f.d))
	} else {
		return fmt.Sprint(f.n)
	}
}

func (f *SdigitsFloat) Description() string {
	s := fmt.Sprintf("%.*G", int(f.sd), f.n)
	return s
}

func (f *Fraction) FractionValue() *Fraction {
	return f
}

func (f *SdigitsFloat) FractionValue() *Fraction {
	return nil
}

func (f *Fraction) FloatValue() *SdigitsFloat {
	sdf := new(SdigitsFloat)
	sdf.n = float64(f.n) / float64(f.d)
	sdf.sd = 15
	return sdf
}

func (f *SdigitsFloat) FloatValue() *SdigitsFloat {
	return f
}

//Math
func MakeFraction(numerator, denominator int64) *Fraction {
	if denominator == 0 {
		fmt.Println("Can't create fraction with denominator=0")
		return nil
	}
	if denominator == 1 {
		nfrac := new(Fraction)
		nfrac.n = numerator
		nfrac.d = 1
		nfrac.om = 1
		return nfrac
	}
	nf := primeFactors(numerator)
	df := primeFactors(denominator)
	var cf int64 = 1

	ni := 0
	di := 0

	for {
		if nf[ni] < df[di] {
			if len(nf) == ni+1 {
				break
			}
			ni++
		} else if nf[ni] > df[di] {
			if len(df) == di+1 {
				break
			}
			di++
		} else {
			cf *= nf[ni]
			if len(nf) == ni+1 {
				break
			}
			if len(df) == di+1 {
				break
			}
			ni++
			di++
		}
	}
	if denominator < 0 && numerator > 0 {
		cf = cf * -1
	}
	nfrac := new(Fraction)
	nfrac.n = numerator / cf
	nfrac.d = denominator / cf
	nfrac.om = cf
	return nfrac
}

func MakeSDFloat(n float64, sd uint8) *SdigitsFloat {
	nsdf := new(SdigitsFloat)
	nsdf.n = n
	nsdf.sd = sd
	return nsdf
}

func add(args ...number) number {
	a, b := args[0], args[1]
	if a == nil || b == nil {
		return nil
	}
	if an, bn := a.FractionValue(), b.FractionValue(); an != nil && bn != nil {
		am, bm := getFactorsForCommonDenominator(an, bn)
		return MakeFraction(an.n*am+bn.n*bm, an.d*am)
	}
	sdf := new(SdigitsFloat)
	af, bf := a.FloatValue(), b.FloatValue()
	sdf.n = af.n + bf.n
	sdf.sd = lowestSD(af.sd, bf.sd)
	return sdf
}

func subtract(args ...number) number {
	a, b := args[0], args[1]
	if a == nil || b == nil {
		return nil
	}
	if an, bn := a.FractionValue(), b.FractionValue(); an != nil && bn != nil {
		am, bm := getFactorsForCommonDenominator(an, bn)
		return MakeFraction(an.n*am-bn.n*bm, an.d*am)
	}
	sdf := new(SdigitsFloat)
	af, bf := a.FloatValue(), b.FloatValue()
	sdf.n = af.n - bf.n
	sdf.sd = lowestSD(af.sd, bf.sd)
	return sdf
}

func multiply(args ...number) number {
	a, b := args[0], args[1]
	if a == nil || b == nil {
		return nil
	}
	if an, bn := a.FractionValue(), b.FractionValue(); an != nil && bn != nil {
		return MakeFraction(an.n*bn.n, an.d*bn.d)
	}
	sdf := new(SdigitsFloat)
	af, bf := a.FloatValue(), b.FloatValue()
	sdf.n = af.n * bf.n
	sdf.sd = lowestSD(af.sd, bf.sd)
	return sdf
}

func divide(args ...number) number {
	a, b := args[0], args[1]
	if a == nil || b == nil {
		return nil
	}
	if an, bn := a.FractionValue(), b.FractionValue(); an != nil && bn != nil {
		return MakeFraction(an.n*bn.d, an.d*bn.n)
	}
	sdf := new(SdigitsFloat)
	af, bf := a.FloatValue(), b.FloatValue()
	sdf.n = af.n / bf.n
	sdf.sd = lowestSD(af.sd, bf.sd)
	return sdf
}

func raiseToPower(args ...number) number {
	a, b := args[0], args[1]
	if a == nil || b == nil {
		return nil
	}
	sdf := new(SdigitsFloat)
	af, bf := a.FloatValue(), b.FloatValue()
	sdf.n = math.Pow(af.n, bf.n)
	sdf.sd = lowestSD(af.sd, bf.sd)
	return sdf
}

func ln(args ...number) number {
	a := args[0]
	if a == nil {
		return nil
	}
	sdf := new(SdigitsFloat)
	af := a.FloatValue()
	sdf.n = math.Log(af.n)
	sdf.sd = af.sd
	return sdf
}

func lg(args ...number) number {
	a := args[0]
	if a == nil {
		return nil
	}
	sdf := new(SdigitsFloat)
	af := a.FloatValue()
	sdf.n = math.Log10(af.n)
	sdf.sd = af.sd
	return sdf
}

//Utillity functions
func getFactorsForCommonDenominator(a, b *Fraction) (am, bm int64) {
	am, bm = 1, 1
	if a.d != b.d {
		var af, bf []int64 = primeFactors(a.d), primeFactors(b.d)
		ad, bd := false, false
		ai, bi := 0, 0
		for {
			if af[ai] == bf[bi] && !(ad || bd) {
				ai++
				bi++
			} else if (af[ai] < bf[bi] || bd) && !ad {
				bm = bm * af[ai]
				ai++
			} else if (bf[bi] < af[ai] || ad) && !bd {
				am = am * bf[bi]
				bi++
			}

			if ai == len(af) {
				ai--
				ad = true
			}
			if bi == len(bf) {
				bi--
				bd = true
			}
			if ad && bd {
				break
			}
		}
	}
	return
}

func lowestSD(a, b uint8) uint8 {
	if a <= b {
		return a
	} else {
		return b
	}
}

//Prime functions
func calcNextPrime() {
	if cap(primes) == 0 {
		primes = make([]int64, 1)
		primes[0] = 2
	}

	var i int64
	i = primes[len(primes)-1] + 1
	for {
		isPrime := true
		for c := 0; float64(primes[c]) <= math.Sqrt(float64(i)); c++ {
			if i%primes[c] == 0 {
				isPrime = false
				break
			}

		}

		if isPrime {
			primes = append(primes, i)
			return
		}

		i++
	}
}

func primeFactors(n int64) []int64 {
	if n == 0 {
		a := make([]int64, 1)
		a[0] = 0
		return a
	} else if n == 1 {
		a := make([]int64, 1)
		a[0] = 1
		return a
	}

	ret := make([]int64, 0, 1)
	if cap(primes) == 0 {
		calcNextPrime()
	}
	if n < 0 {
		ret = append(ret, -1)
		n = -n
	}

	i := 0

	for float64(primes[i]) <= math.Sqrt(float64(n)) {
		if n%primes[i] == 0 {
			n = n / primes[i]
			ret = append(ret, primes[i])
			i = 0
		} else {
			if i+1 == len(primes) {
				calcNextPrime()
			}
			i++
		}
	}
	if n > 1 {
		ret = append(ret, n)
	}
	return ret
}
