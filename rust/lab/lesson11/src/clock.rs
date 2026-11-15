use std::time::{SystemTime, UNIX_EPOCH};

/// 현재 시각을 RFC 3339 UTC 문자열로 돌려줍니다. 소수점 아래는 0이 아닐 때만, 끝의 0을 뺀 채 붙입니다.
pub fn now_rfc3339() -> String {
    let elapsed = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap_or_default();
    format_timestamp(elapsed.as_secs(), elapsed.subsec_nanos())
}

fn format_timestamp(seconds: u64, nanos: u32) -> String {
    let days = (seconds / 86_400) as i64;
    let rest = seconds % 86_400;
    let (year, month, day) = civil_from_days(days);
    let mut text = format!(
        "{year:04}-{month:02}-{day:02}T{:02}:{:02}:{:02}",
        rest / 3_600,
        rest % 3_600 / 60,
        rest % 60
    );
    if nanos != 0 {
        let fraction = format!("{nanos:09}");
        text.push('.');
        text.push_str(fraction.trim_end_matches('0'));
    }
    text.push('Z');
    text
}

/// 1970-01-01부터 지난 날짜 수를 (년, 월, 일)로 바꿉니다. (Howard Hinnant의 civil_from_days 알고리즘)
fn civil_from_days(days: i64) -> (i64, u32, u32) {
    let z = days + 719_468;
    let era = z.div_euclid(146_097);
    let day_of_era = z.rem_euclid(146_097);
    let year_of_era =
        (day_of_era - day_of_era / 1_460 + day_of_era / 36_524 - day_of_era / 146_096) / 365;
    let day_of_year = day_of_era - (365 * year_of_era + year_of_era / 4 - year_of_era / 100);
    let month_index = (5 * day_of_year + 2) / 153;
    let day = (day_of_year - (153 * month_index + 2) / 5 + 1) as u32;
    let month = (if month_index < 10 {
        month_index + 3
    } else {
        month_index - 9
    }) as u32;
    let year = year_of_era + era * 400 + i64::from(month <= 2);
    (year, month, day)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn formats_known_timestamps() {
        assert_eq!(format_timestamp(0, 0), "1970-01-01T00:00:00Z");
        assert_eq!(format_timestamp(951_782_400, 0), "2000-02-29T00:00:00Z");
        assert_eq!(format_timestamp(1_000_000_000, 0), "2001-09-09T01:46:40Z");
        assert_eq!(format_timestamp(1_700_000_000, 0), "2023-11-14T22:13:20Z");
        assert_eq!(format_timestamp(4_102_444_800, 0), "2100-01-01T00:00:00Z");
    }

    #[test]
    fn fraction_is_trimmed_like_go() {
        assert_eq!(
            format_timestamp(1_700_000_000, 123_456_700),
            "2023-11-14T22:13:20.1234567Z"
        );
        assert_eq!(
            format_timestamp(1_700_000_000, 500_000_000),
            "2023-11-14T22:13:20.5Z"
        );
        assert_eq!(
            format_timestamp(1_700_000_000, 1),
            "2023-11-14T22:13:20.000000001Z"
        );
    }

    #[test]
    fn now_has_rfc3339_shape() {
        let now = now_rfc3339();
        assert!(now.starts_with("20"));
        assert!(now.ends_with('Z'));
        assert_eq!(&now[10..11], "T");
    }
}
