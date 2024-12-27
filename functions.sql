--Get products for merchant
CREATE OR REPLACE FUNCTION getProductsMerchant_fn(userId uuid,brandName text,categoryName text)
RETURNS TABLE (product_id uuid,category_id uuid,brand_id uuid,user_id uuid,product_name text,price numeric,rating numeric) AS
$$
BEGIN
	
	IF brandName != '' AND categoryName != '' THEN
		RETURN QUERY	
		SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id)
		WHERE category_name = categoryName AND brand_name = brandName AND p.user_id = userId;
	ELSEIF brandName = '' AND categoryName != '' THEN
		RETURN QUERY	
	    SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id)
		WHERE category_name = categoryName AND p.user_id = userId ;
	ELSEIF brandName != '' AND categoryName = '' THEN
		RETURN QUERY	
	    SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id)
		WHERE brand_Name = brandName AND p.user_id = userId ;
	ELSE 
		RETURN QUERY	
		SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id)
		WHERE p.user_id = userId;
	END IF;	
END;
$$
LANGUAGE plpgsql;



--Filter products for User
CREATE OR REPLACE FUNCTION filterProductsUser_fn(requiredPrice numeric,requiredRating numeric)
RETURNS TABLE (product_id uuid,category_id uuid,brand_id uuid,user_id uuid,product_name text,price numeric,rating numeric) AS
$$
BEGIN

	IF requiredPrice = 0 AND requiredRating != 0 THEN
		RETURN QUERY	
	    SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id)
		WHERE p.rating >= requiredRating;
	ELSEIF requiredPrice != 0 AND requiredRating = 0 THEN
		RETURN QUERY	
	    SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id)
		WHERE p.price >= requiredPrice;
	ELSE 
		RETURN QUERY	
		SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id)
		WHERE p.price >= requiredPrice AND p.rating >= requiredRating;
	END IF;	
END;
$$
LANGUAGE plpgsql;



--Get products for User
CREATE OR REPLACE FUNCTION getProductsUser_fn(brandName text,categoryName text)
RETURNS TABLE (product_id uuid,category_id uuid,brand_id uuid,user_id uuid,product_name text,price numeric,rating numeric) AS
$$
BEGIN
	
	IF brandName != '' AND categoryName != '' THEN
		RETURN QUERY	
		SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id)
		WHERE category_name = categoryName AND brand_name = brandName;
	ELSEIF brandName = '' AND categoryName != '' THEN
		RETURN QUERY	
	    SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id)
		WHERE category_name = categoryName;
	ELSEIF brandName != '' AND categoryName = '' THEN
		RETURN QUERY	
	    SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id)
		WHERE brand_Name = brandName;
	ELSE 
		RETURN QUERY	
		SELECT p.product_id,p.category_id,p.brand_id,p.user_id,p.product_name,p.price,p.rating FROM products AS p 
	    INNER JOIN categories USING(category_id) 
		INNER JOIN brands USING(Brand_id);
	END IF;	
END;
$$
LANGUAGE plpgsql;
