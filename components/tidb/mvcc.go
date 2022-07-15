package tidb

//type mvcc struct {
//	Key      string `json:"key"`
//	RegionID int    `json:"region_id"`
//	Value    struct {
//		Info struct {
//			Writes []struct {
//				StartTs    int    `json:"start_ts"`
//				CommitTs   int    `json:"commit_ts"`
//				ShortValue string `json:"short_value"`
//			} `json:"writes"`
//		} `json:"info"`
//	} `json:"value"`
//}
//
//func (r *Runner) tidbTableMVCC(tbl string) error {
//	validate := func(input string) error {
//		_, err := strconv.ParseFloat(input, 64)
//		return err
//	}
//
//	templates := &promptui.PromptTemplates{
//		Prompt:  "{{ . }} ",
//		Valid:   "{{ . | green }} ",
//		Invalid: "{{ . | red }} ",
//		Success: "{{ . | bold }} ",
//	}
//
//	prompt := promptui.Prompt{
//		Label:     "which?",
//		Templates: templates,
//		Validate:  validate,
//	}
//
//	key, err := prompt.Run()
//	if err != nil {
//		return err
//	}
//
//	if err := r.displayTiDBKeyMVCC(tbl, key); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (r *Runner) displayTiDBKeyMVCC(tb string, key string) error {
//	req := fmt.Sprintf("http://%s:%d/mvcc/key/%s/%s", strings.Split(r.TidbCluster[0].Host, ":")[0], r.TidbCluster[0].StatusPort, strings.ReplaceAll(tb, ".", "/"), key)
//	body, err := util.GetHttp(req)
//	if err != nil {
//		return err
//	}
//	var cc mvcc
//	err = json.Unmarshal(body, &cc)
//
//	ms := cc.Value.Info.Writes
//
//	templates := &promptui.SelectTemplates{
//		Label:    "which?",
//		Active:   "\U0001F550  {{ .StartTs | red }} {{ .CommitTs | red }}",
//		Inactive: "{{ .StartTs | white }} {{ .CommitTs | white }}",
//		Selected: "\U0001F550  {{ .StartTs | red | white }} {{ .CommitTs | red | white }}",
//		Details:  ``,
//	}
//
//	searcher := func(input string, index int) bool {
//		MVCC := ms[index]
//		name := strings.Replace(strings.ToLower(strconv.FormatInt(int64(MVCC.StartTs), 10)), " ", "", -1)
//		input = strings.Replace(strings.ToLower(input), " ", "", -1)
//		return strings.Contains(name, input)
//	}
//
//	prompt := promptui.Select{
//		Label:     "MVCC",
//		Items:     ms,
//		Templates: templates,
//		Size:      4,
//		Searcher:  searcher,
//	}
//
//	i, _, err := prompt.Run()
//
//	l := list.NewWriter()
//	l.SetStyle(list.StyleBulletCircle)
//
//	var s []interface{}
//	s = append(s, fmt.Sprintf("startTs:  %s", oracle.GetTimeFromTS(uint64(cc.Value.Info.Writes[i].StartTs))))
//	s = append(s, fmt.Sprintf("commitTs: %s", oracle.GetTimeFromTS(uint64(cc.Value.Info.Writes[i].CommitTs))))
//	s = append(s, "snapshot")
//	l.AppendItems(s)
//	shortValue, err := r.decodeKeyMVCCShortValue(tb, cc.Value.Info.Writes[i].ShortValue)
//	l.Indent()
//	var shot []interface{}
//	for k, v := range shortValue {
//		shot = append(shot, fmt.Sprintf("%s: %s", k, v))
//	}
//	l.AppendItems(shot)
//	fmt.Println(l.Render())
//
//	p := promptui.Prompt{
//		Label:     "return",
//		IsConfirm: true,
//	}
//	result, _ := p.Run()
//	if result == "y" {
//		if err := r.displayTiDBKeyMVCC(tb, key); err != nil {
//			return err
//		}
//	} else {
//		sys.Exit()
//	}
//
//	return nil
//}
//
//func (r *Runner) decodeKeyMVCCShortValue(tblName string, shortValue string) (map[string]string, error) {
//	tblInfo, err := TableInfo(tblName, r.TidbCluster[0].Host, r.TidbCluster[0].StatusPort)
//	if err != nil {
//		return nil, err
//	}
//	m, err := decodeMVCC(tblInfo, shortValue)
//	if err != nil {
//		return nil, err
//	}
//	return m, nil
//}
//
//func decodeMVCC(tblInfo *model.TableInfo, id string) (map[string]string, error) {
//	bs, err := base64.StdEncoding.DecodeString(id)
//	if err != nil {
//		return nil, err
//	}
//	colMap := make(map[int64]*types.FieldType, 3)
//	for _, col := range tblInfo.Columns {
//		colMap[col.ID] = &col.FieldType
//	}
//	r, err := decodeRow(bs, colMap, oracle.UTC)
//	m := make(map[string]string)
//	for _, col := range tblInfo.Columns {
//		if v, ok := r[col.ID]; ok {
//			if v.IsNull() {
//				m[col.Name.L] = "NULL"
//				continue
//			}
//			ss, err := v.ToString()
//			if err != nil {
//				m[col.Name.L] = "error: " + err.Error() + fmt.Sprintf("datum: %#v", v)
//				continue
//			}
//			m[col.Name.L] = ss
//		}
//	}
//	return m, nil
//}
//
//func decodeRow(b []byte, cols map[int64]*types.FieldType, loc *oracle.Location) (map[int64]types.Datum, error) {
//	if !rowcodec.IsNewFormat(b) {
//		return tablecodec.DecodeRowWithMap(b, cols, loc, nil)
//	}
//	return tablecodec.DecodeRowWithMapNew(b, cols, loc, nil)
//}
