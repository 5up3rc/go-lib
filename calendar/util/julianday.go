/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package util

import (
	"math"
	"time"
)

// GetDateFromJulianDay 从儒略日中获取公历的日期
func GetDateFromJulianDay(jd float64) (yy, mm, dd int) {
	/*
	 * This algorithm is taken from
	 * "Numerical Recipes in c, 2nd Ed." (1992), pp. 14-15
	 * and converted to integer math.
	 * The electronic version of the book is freely available
	 * at http://www.nr.com/ , under "Obsolete Versions - Older
	 * book and code versions.
	 */
	const JD_GREG_CAL = 2299161
	const JB_MAX_WITHOUT_OVERFLOW = 107374182
	julian := int64(math.Floor(jd + 0.5))

	var ta, jalpha, tb, tc, td, te int64

	if julian >= JD_GREG_CAL {
		jalpha = (4*(julian-1867216) - 1) / 146097
		ta = int64(julian) + 1 + jalpha - jalpha/4
	} else if julian < 0 {
		ta = julian + 36525*(1-julian/36525)
	} else {
		ta = julian
	}
	tb = ta + 1524
	if tb <= JB_MAX_WITHOUT_OVERFLOW {
		tc = (tb*20 - 2442) / 7305
	} else {
		tc = int64((uint64(tb)*20 - 2442) / 7305)
	}
	td = 365*tc + tc/4
	te = ((tb - td) * 10000) / 306001

	dd = int(tb - td - (306001*te)/10000)

	mm = int(te - 1)
	if mm > 12 {
		mm -= 12
	}
	yy = int(tc - 4715)
	if mm > 2 {
		yy--
	}
	if julian < 0 {
		yy -= int(100 * (1 - julian/36525))
	}

	return
}

// GetTimeFromJulianDay 从儒略日中获取时间 时分秒
func GetTimeFromJulianDay(jd float64) (hour, minute, second int) {
	frac := jd - math.Floor(jd)
	s := int(math.Floor(frac * 24.0 * 60.0 * 60.0))

	hour = ((s / (60 * 60)) + 12) % 24
	minute = (s / (60)) % 60
	second = s % 60
	return
}

// GetDateTimeFromJulianDay 将儒略日转换为 time.Time
// 其中包含了 TT 到 UTC 的转换
func GetDateTimeFromJulianDay(jd float64) time.Time {
	yy, mm, dd := GetDateFromJulianDay(jd)
	//  TT -> UTC
	jd -= GetDeltaT(yy, mm) / 86400
	yy, mm, dd = GetDateFromJulianDay(jd)
	h, m, s := GetTimeFromJulianDay(jd)
	return time.Date(yy, time.Month(mm), dd, h, m, s, 0, time.UTC)
}
