package model

import (
	"database/sql"
	"log"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

type CasbinModel struct {
	DB *sql.DB
}

type CasbinPolicy struct {
	Sub, Obj, Act string
}

func (a *CasbinModel) Init() (*casbin.Enforcer, error) {
	tableName := "casbin"
	adpt, err := sqladapter.NewAdapter(a.DB, "sqlite3", tableName)
	if err != nil {
		panic(err)
	}

	m, err := model.NewModelFromString(`
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow)) && !some(where (p.eft == deny))
	
	[matchers]
    m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*") || r.sub=="admin"
	`)
	if err != nil {
		log.Fatalf("error: model: %s", err)
	}

	enforcer, err := casbin.NewEnforcer(m, adpt)
	if err != nil {
		log.Fatalf("New enforcer error: %s", err)
	}

	return enforcer, nil
}

// GetPolicies get all policies by role name
func (a *CasbinModel) GetPolicies(roleName string) []*CasbinPolicy {
	enforer, err := a.Init()
	if err != nil {
		log.Fatal("casbin init error")
	}

	var rules [][]string

	if roleName == "" {
		rules = enforer.GetPolicy()
	} else {
		rules = enforer.GetFilteredPolicy(0, roleName)
	}

	policies := []*CasbinPolicy{}
	for _, rule := range rules {
		policy := &CasbinPolicy{Sub: rule[0], Obj: rule[1], Act: rule[2]}
		policies = append(policies, policy)
	}

	return policies
}

// GetPoliciesOrderBy order by role name from database
func (m *CasbinModel) GetPoliciesOrderBy(roleName string) []*CasbinPolicy {
	var q string

	if roleName == "Administrator" || roleName == "" {
		q = `SELECT v0,v1,v2 FROM casbin WHERE p_type='p'
				 ORDER BY v0,v1`
	} else {
		q = `SELECT v0,v1,v2 FROM casbin 
				 WHERE p_type='p' AND v0='` + roleName + `'
				 ORDER BY v0,v1`
	}

	rows, err := m.DB.Query(q)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer rows.Close()

	policies := []*CasbinPolicy{}
	for rows.Next() {
		p := &CasbinPolicy{}
		if err := rows.Scan(&p.Sub, &p.Obj, &p.Act); err != nil {
			log.Println(err)
			return nil
		}
		policies = append(policies, p)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil
	}

	return policies
}

func (m *CasbinModel) Edit(old, new []string) error {
	enforcer, err := m.Init()
	if err != nil {
		log.Fatal("casbin init error")
	}

	// UpdatePolicy not working,so delete then insert
	//updated, err := enforcer.UpdatePolicy(old, new)

	_, err = enforcer.RemovePolicy(old)
	if err != nil {
		log.Println(err)
	}
	_, err = enforcer.AddPolicy(new)
	if err != nil {
		log.Println(err)
	}

	return err
}
