package pagerduty

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyEventRuleCreate,
		Read:   resourcePagerDutyEventRuleRead,
		Update: resourcePagerDutyEventRuleUpdate,
		Delete: resourcePagerDutyEventRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"action_json": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition_json": {
				Type:     schema.TypeString,
				Required: true,
			},
			"advanced_condition_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"catch_all": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func buildEventRuleStruct(d *schema.ResourceData) *pagerduty.EventRule {
	eventRule := &pagerduty.EventRule{
		Actions:   expandString(d.Get("action_json").(string)),
		Condition: expandString(d.Get("condition_json").(string)),
	}

	if attr, ok := d.GetOk("advanced_condition"); ok {
		eventRule.AdvancedCondition = expandString(attr.(string))
	}

	if attr, ok := d.GetOk("catch_all"); ok {
		eventRule.CatchAll = attr.(bool)
	}

	return eventRule
}

func expandString(v string) []interface{} {
	var obj []interface{}
	err := json.Unmarshal([]byte(v), &obj)

	if err != nil {
		log.Printf(string(err.Error()))
	}

	return obj
}

func flattenSlice(v []interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf(string(err.Error()))
	}
	return string(b)
}

func resourcePagerDutyEventRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	eventRule := buildEventRuleStruct(d)

	log.Printf("[INFO] Creating PagerDuty event rule: %s", "eventRule")

	eventRule, _, err := client.EventRules.Create(eventRule)
	if err != nil {
		return err
	}

	d.SetId(eventRule.ID)

	return resourcePagerDutyEventRuleRead(d, meta)
}

func resourcePagerDutyEventRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty event rule: %s", d.Id())

	var eventRule *pagerduty.EventRule

	resp, _, err := client.EventRules.List()
	if err != nil {
		return err
	}
	for _, rule := range resp.EventRules {
		if rule.ID == d.Id() {
			d.Set("action_json", flattenSlice(rule.Actions))
			d.Set("condition_json", flattenSlice(rule.Condition))
			d.Set("action_json", flattenSlice(rule.Actions))
			d.Set("catch_all", rule.CatchAll)

		}
	}
	if eventRule.ID != d.Id() {
		return handleNotFoundError(err, d)
	}

	return nil
}
func resourcePagerDutyEventRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	eventRule := buildEventRuleStruct(d)

	log.Printf("[INFO] Updating PagerDuty event rule: %s", d.Id())

	if _, _, err := client.EventRules.Update(d.Id(), eventRule); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyEventRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty event rule: %s", d.Id())

	if _, err := client.EventRules.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
