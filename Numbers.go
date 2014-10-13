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
	fractionValue() *Fraction //returns nil when called on SdigitsFloat
	floatValue() *SdigitsFloat
	description() string
}

//Number interface functions
func (f *Fraction) description() string {
	if f.d != 1 && f.om != 1 {
		return fmt.Sprintf("%d/%d = %d/%d = %f", f.n*f.om, f.d*f.om, f.n, f.d, float64(f.n)/float64(f.d))
	} else if f.d != 1 {
		return fmt.Sprintf("%d/%d = %f", f.n, f.d, float64(f.n)/float64(f.d))
	} else {
		return fmt.Sprint(f.n)
	}
}

func (f *SdigitsFloat) description() string {
	s := fmt.Sprintf("%.*G", int(f.sd), f.n)
	return s
}

func (f *Fraction) fractionValue() *Fraction {
	return f
}

func (f *SdigitsFloat) fractionValue() *Fraction {
	return nil
}

func (f *Fraction) floatValue() *SdigitsFloat {
	sdf := new(SdigitsFloat)
	sdf.n = float64(f.n) / float64(f.d)
	sdf.sd = 15
	return sdf
}

func (f *SdigitsFloat) floatValue() *SdigitsFloat {
	return f
}

//Math
func makeFraction(numerator, denominator int64) *Fraction {
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

func add(args ...number) number {
	a := args[0]
	b := args[1]
	if an, bn := a.fractionValue(), b.fractionValue(); an != nil && bn != nil {
		var am, bm int64 = 1, 1

		if an.d != bn.d {
			var af, bf []int64 = primeFactors(an.d), primeFactors(bn.d)

			ad, bd := false, false
			ai, bi := 0, 0
			for {
				if af[ai] == bf[bi] && !(ad || bd) {
					ai++
					bi++
				} else if af[ai] < bf[bi] || bd {
					bm = bm * af[ai]
					ai++
				} else if bf[bi] < af[ai] || ad {
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
		return makeFraction(an.n*am+bn.n*bm, an.d*am)
	}
	sdf := new(SdigitsFloat)
	af, bf := a.floatValue(), b.floatValue()
	sdf.n = af.n + bf.n
	if af.sd <= bf.sd {
		sdf.sd = af.sd
	} else {
		sdf.sd = bf.sd
	}
	return sdf
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
