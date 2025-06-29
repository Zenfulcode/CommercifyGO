-- Revert currency exchange rates to the previous incorrect values
-- This is just for rollback purposes - the old rates were incorrect

UPDATE currencies SET exchange_rate = 0.15 WHERE code = 'DKK';  -- Previous incorrect rate
UPDATE currencies SET exchange_rate = 0.85 WHERE code = 'EUR';  -- Previous rate
UPDATE currencies SET exchange_rate = 0.75 WHERE code = 'GBP';  -- Previous rate
UPDATE currencies SET exchange_rate = 110.0 WHERE code = 'JPY'; -- Previous rate
UPDATE currencies SET exchange_rate = 1.25 WHERE code = 'CAD';  -- Previous rate