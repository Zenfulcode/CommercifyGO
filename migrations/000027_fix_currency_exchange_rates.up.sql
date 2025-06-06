-- Fix incorrect currency exchange rates
-- The previous migration had incorrect exchange rates that were way too low
-- This migration updates them to more realistic values as of June 2025

UPDATE currencies SET exchange_rate = 6.54 WHERE code = 'DKK';  -- Danish Krone: 1 USD = ~6.54 DKK
UPDATE currencies SET exchange_rate = 0.92 WHERE code = 'EUR';  -- Euro: 1 USD = ~0.92 EUR  
UPDATE currencies SET exchange_rate = 0.79 WHERE code = 'GBP';  -- British Pound: 1 USD = ~0.79 GBP
UPDATE currencies SET exchange_rate = 149.50 WHERE code = 'JPY'; -- Japanese Yen: 1 USD = ~149.50 JPY
UPDATE currencies SET exchange_rate = 1.37 WHERE code = 'CAD';  -- Canadian Dollar: 1 USD = ~1.37 CAD