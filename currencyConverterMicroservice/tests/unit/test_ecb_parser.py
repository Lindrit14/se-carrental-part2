"""Parser unit tests against a static ECB-shaped XML fixture."""
from __future__ import annotations

from datetime import date
from decimal import Decimal

from app.infrastructure.ecb_provider import parse_ecb_xml

ECB_XML = """<?xml version="1.0" encoding="UTF-8"?>
<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01"
                 xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
  <gesmes:subject>Reference rates</gesmes:subject>
  <Cube>
    <Cube time="2026-04-28">
      <Cube currency="USD" rate="1.0823"/>
      <Cube currency="GBP" rate="0.8634"/>
      <Cube currency="CHF" rate="0.9612"/>
    </Cube>
  </Cube>
</gesmes:Envelope>
"""


def test_parse_ecb_xml_extracts_date_and_rates() -> None:
    rate_set = parse_ecb_xml(ECB_XML)
    assert rate_set.rate_date == date(2026, 4, 28)
    assert rate_set.rates["USD"] == Decimal("1.0823")
    assert rate_set.rates["CHF"] == Decimal("0.9612")
    assert "EUR" not in rate_set.rates  # EUR is the base, not in feed
