-- Candelize temperature by minute
SELECT 
	DATE_PART('year', dt) as year, 
	DATE_PART('month', dt) as month, 
	DATE_PART('day', dt) as day, 
	DATE_PART('hour', dt) as hour, 
	DATE_PART('minute', dt) as minute, 
	MIN(message::float)::numeric(10,2) AS min_temp,
	AVG(message::float)::numeric(10,2) AS avg_temp,
	MAX(message::float)::numeric(10,2) AS max_temp
FROM rawdata 
WHERE 
	topic = 'croco/cave/temperature' AND
	date_trunc('day', dt) = '2022-01-11'
GROUP BY year, month, day, hour, minute
ORDER BY year, month, day, hour, minute;


