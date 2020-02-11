//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package collectors

import (
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"time"
)

type bucketInfoCollector struct {
	m MetaCollector

	basicstatsDataused         *prometheus.Desc
	basicstatsDiskfetches      *prometheus.Desc
	basicstatsDiskused         *prometheus.Desc
	basicstatsItemcount        *prometheus.Desc
	basicstatsMemused          *prometheus.Desc
	basicstatsOpspersec        *prometheus.Desc
	basicstatsQuotapercentused *prometheus.Desc
}

func (c *bucketInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration

	ch <- c.basicstatsDataused
	ch <- c.basicstatsDiskfetches
	ch <- c.basicstatsDiskused
	ch <- c.basicstatsItemcount
	ch <- c.basicstatsMemused
	ch <- c.basicstatsOpspersec
	ch <- c.basicstatsQuotapercentused
}

func (c *bucketInfoCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()
	log.Info("Collecting buckets metrics...")

	buckets, err := c.m.client.Buckets()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0)
		log.With("error", err).Error("failed to scrape buckets")
		return
	}

	for _, bucket := range buckets {
		log.Debugf("Collecting %s bucket metrics...", bucket.Name)

		ch <- prometheus.MustNewConstMetric(c.basicstatsDataused, prometheus.GaugeValue, bucket.BucketBasicStats.DataUsed, bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.basicstatsDiskfetches, prometheus.GaugeValue, bucket.BucketBasicStats.DiskFetches, bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.basicstatsDiskused, prometheus.GaugeValue, bucket.BucketBasicStats.DiskUsed, bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.basicstatsItemcount, prometheus.GaugeValue, bucket.BucketBasicStats.ItemCount, bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.basicstatsMemused, prometheus.GaugeValue, bucket.BucketBasicStats.MemUsed, bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.basicstatsOpspersec, prometheus.GaugeValue, bucket.BucketBasicStats.OpsPerSec, bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.basicstatsQuotapercentused, prometheus.GaugeValue, bucket.BucketBasicStats.QuotaPercentUsed, bucket.Name)
	}

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())

}

func NewBucketInfoCollector(client util.Client) prometheus.Collector {
	const subsystem = "bucketinfo"

	return &bucketInfoCollector{
		m: MetaCollector{
			client: client,
			up: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE + subsystem, "", "up"),
				"Couchbase buckets API is responding",
				nil,
				nil,
			),
			scrapeDuration: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE + subsystem, "", "scrape_duration_seconds"),
				"Scrape duration in seconds",
				nil,
				nil,
			),
		},
		basicstatsDataused: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE + subsystem, "", "basic_dataused_bytes"),
			"basic_dataused",
			[]string{"bucket"},
			nil,
		),
		basicstatsDiskfetches: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE + subsystem, "", "basic_diskfetches"),
			"basic_diskfetches",
			[]string{"bucket"},
			nil,
		),
		basicstatsDiskused: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE + subsystem, "", "basic_diskused_bytes"),
			"basic_diskused",
			[]string{"bucket"},
			nil,
		),
		basicstatsItemcount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE + subsystem, "", "basic_itemcount"),
			"basic_itemcount",
			[]string{"bucket"},
			nil,
		),
		basicstatsMemused: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE + subsystem, "", "basic_memused_bytes"),
			"basic_memused",
			[]string{"bucket"},
			nil,
		),
		basicstatsOpspersec: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE + subsystem, "", "basic_opspersec"),
			"basic_opspersec",
			[]string{"bucket"},
			nil,
		),
		basicstatsQuotapercentused: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE + subsystem, "", "basic_quota_user_percent"),
			"basic_quotapercentused",
			[]string{"bucket"},
			nil,
		),
	}
}
