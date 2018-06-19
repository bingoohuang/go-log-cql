CREATE KEYSPACE blackcat with replication =
    {'class': 'SimpleStrategy', 'replication_factor': 1 };

USE blackcat;

-- 日志中的异常
DROP TABLE IF EXISTS event_log_exception;
CREATE TABLE event_log_exception (
    logId text,
    hostname text,
    logger text,
    tcode text,
    tid text,
    exceptionNames text,
    contextLogs text,
    timestamp text,
    PRIMARY KEY(logId)
) WITH default_time_to_live = 2592000;


INSERT INTO event_log_exception(logId, logger, tcode, tid, exceptionNames, contextLogs, timestamp)
VALUES ('1', 'yoga', '18001', '18001', 'xxxx', 'yyyyy', '2018-06-19 22:21:04');

