package observability

import (
	"time"

	"github.com/azahir21/go-backend-boilerplate/ent"
	"github.com/sirupsen/logrus"
)

// DBStatsCollector periodically collects database connection pool statistics
type DBStatsCollector struct {
	client  *ent.Client
	metrics *Metrics
	log     *logrus.Logger
	stopCh  chan struct{}
}

// NewDBStatsCollector creates a new database stats collector
func NewDBStatsCollector(client *ent.Client, metrics *Metrics, log *logrus.Logger) *DBStatsCollector {
	return &DBStatsCollector{
		client:  client,
		metrics: metrics,
		log:     log,
		stopCh:  make(chan struct{}),
	}
}

// Start begins collecting database statistics periodically
func (c *DBStatsCollector) Start() {
	go c.collectLoop()
}

// Stop stops the statistics collection
func (c *DBStatsCollector) Stop() {
	close(c.stopCh)
}

// collectLoop collects database statistics every 30 seconds
func (c *DBStatsCollector) collectLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.collectStats()
		case <-c.stopCh:
			return
		}
	}
}

// collectStats is a placeholder for collecting database connection pool statistics
// Note: Ent ORM doesn't expose the underlying database/sql.DB directly, making it
// difficult to collect connection pool statistics without modifying the database layer.
// As a workaround, you can:
// 1. Use a custom ent driver that wraps the SQL driver
// 2. Instrument database operations at the query level
// 3. Use database-specific exporters (postgres_exporter, mysql_exporter)
func (c *DBStatsCollector) collectStats() {
	// Currently this is a placeholder. In production, consider using:
	// - Database-specific exporters (postgres_exporter, mysqld_exporter)
	// - Custom ent hooks to track query metrics
	// - Instrumentation at the SQL driver level before ent.Open()
	
	c.log.Debug("Database connection pool stats collection not implemented for ent ORM")
}
