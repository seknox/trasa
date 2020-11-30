package middlewares

//TODO move this test to integration-tests
//Use real database and redis

// func Test_getUserContext(t *testing.T) {
// 	userstore := users.InitStoreMock()
// 	orgstore := orgs.InitStoreMock()
// 	_ = redis.InitStoreMock()
// 	//
// 	mockUser := models.User{
// 		ID:         "someUserID",
// 		OrgID:      "someOrgID",
// 		UserName:   "testUname",
// 		FirstName:  "Bha",
// 		MiddleName: "",
// 		LastName:   "Ach",
// 		Email:      "user@example.com",
// 		UserRole:   "orgAdmin",
// 		Status:     true,
// 		IdpName:    "trasa",
// 	}
// 	mockUser2 := models.User{
// 		ID:         "someUserID2",
// 		OrgID:      "someOrgID2",
// 		UserName:   "testUname",
// 		FirstName:  "Bha",
// 		MiddleName: "",
// 		LastName:   "Ach",
// 		Email:      "user@example.com",
// 		UserRole:   "orgAdmin",
// 		Status:     true,
// 		IdpName:    "trasa",
// 	}
// 	mockOrg := models.Org{
// 		ID:             "someOrgID",
// 		OrgName:        "testOrg",
// 		Domain:         "example.com",
// 		PrimaryContact: "user@example.com",
// 		Timezone:       "Asia/Kathmandu",
// 		PhoneNumber:    "12345678",
// 	}

// 	mockUC := models.UserContext{
// 		User: &mockUser,
// 		Org:  mockOrg,
// 	}
// 	mockUC2 := models.UserContext{
// 		User: &mockUser2,
// 		Org:  mockOrg,
// 	}

// 	userstore.
// 		On("GetFromID", "", "").
// 		Return(&models.User{}, sql.ErrNoRows)

// 	orgstore.
// 		On("Get", "").
// 		Return(&models.Org{}, sql.ErrNoRows)

// 	// var nulUC models.UserContext
// 	// nulUC.User.ID = "shouldfail"

// 	type args struct {
// 		context models.UserContext
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    models.UserContext
// 		wantErr bool
// 	}{
// 		{
// 			name:    "should fail on different user context",
// 			args:    args{mockUC2},
// 			want:    mockUC,
// 			wantErr: true,
// 		},
// 		{
// 			name:    "with valid userID orgID",
// 			args:    args{mockUC},
// 			want:    mockUC,
// 			wantErr: false,
// 		},
// 		// TODO: Add test cases.
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			setSession, setCsrf, err := auth.SetSession(tt.args.context)
// 			fmt.Println("setSession: ", setSession)
// 			if err != nil {
// 				t.Errorf("failed setting session: %v", err)
// 			}

// 			got, err := validateAndGetUserContext(setSession, setCsrf)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("getUserContext() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.wantErr {
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("getUserContext() got = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
