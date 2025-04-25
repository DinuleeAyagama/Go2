// Product handlers
func getProducts(w http.ResponseWriter, r *http.Request) {
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	respondWithJSON(w, http.StatusOK, product)
}

// User handlers
func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Hash password before saving
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password")
		return
	}
	user.Password = hashedPassword

	if err := db.Create(&user).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var user User
	if err := db.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if !checkPasswordHash(credentials.Password, user.Password) {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := generateToken(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

// Cart handlers
func getCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	var cart Cart
	if err := db.Where("user_id = ? AND status = ?", userID, "cart").Preload("Products").First(&cart).Error; err != nil {
		// Create new cart if not exists
		cart = Cart{
			Order: Order{
				UserID: userID,
				Status: "cart",
			},
		}
		db.Create(&cart)
	}

	respondWithJSON(w, http.StatusOK, cart)
}

func addToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	var request struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get or create cart
	var cart Cart
	if err := db.Where("user_id = ? AND status = ?", userID, "cart").First(&cart).Error; err != nil {
		cart = Cart{
			Order: Order{
				UserID: userID,
				Status: "cart",
			},
		}
		db.Create(&cart)
	}

	// Get product
	var product Product
	if err := db.First(&product, request.ProductID).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Add product to cart
	if err := db.Model(&cart).Association("Products").Append(&product); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not add product to cart")
		return
	}

	// Update total price
	cart.TotalPrice += product.Price * float64(request.Quantity)
	db.Save(&cart)

	respondWithJSON(w, http.StatusOK, cart)
}