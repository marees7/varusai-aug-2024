
--Get products for merchant
CREATE OR REPLACE FUNCTION getProducts_fn(
    userId uuid,
    brandName text,
    categoryName text,
    targetPrice numeric,
    targetRating numeric,
    targetLimit int,
    targetOffset int
)
RETURNS TABLE (
    product_id uuid,
    category_id uuid,
    brand_id uuid,
    user_id uuid,
    product_name text,
    price numeric,
    rating numeric
) AS $$
BEGIN
	If userId IS NOT NULL Then
		 RETURN QUERY 
  		  SELECT 
       		 p.product_id,
       		 p.category_id,
        	 p.brand_id,
       		 p.user_id,
       		 p.product_name,
       		 p.price,
       		 p.rating 
    		FROM products AS p 
    		INNER JOIN categories USING(category_id) 
   			INNER JOIN brands USING(brand_id)
   			WHERE 
			(
				p.user_id = userId
    			AND (brandName = '' OR brand_name = brandName)
    			AND (categoryName = '' OR category_name = categoryName)
    			AND (targetPrice = 0 OR p.price >= targetPrice)
    			AND (targetRating = 0 OR p.rating >= targetRating)
			)
			LIMIT targetLimit 
   		    OFFSET targetOffset;
	ElseIf (userId IS NULL AND brandName = '' AND categoryName = '' AND targetPrice = 0 AND targetRating = 0)  Then
	 		RETURN QUERY 
  	 		 SELECT 
      		    p.product_id,
     		    p.category_id,
       			p.brand_id,
     		   	p.user_id,
      		  	p.product_name,
      		  	p.price,
      		  	p.rating 
  			  FROM products AS p 
    		  INNER JOIN categories USING(category_id) 
  			  INNER JOIN brands USING(brand_id)
			  WHERE p.product_name != ''
			  LIMIT targetLimit 
              OFFSET targetOffset;
	Else 
	 		RETURN QUERY 
  	  		SELECT 
        	p.product_id,
        	p.category_id,
        	p.brand_id,
        	p.user_id,
        	p.product_name,
        	p.price,
        	p.rating 
    		FROM products AS p 
    		INNER JOIN categories USING(category_id) 
    		INNER JOIN brands USING(brand_id)
    		WHERE 
			(
				userId IS NULL
    			AND (brandName = '' OR brand_name = brandName)
    			AND (categoryName = '' OR category_name = categoryName)
    			AND (targetPrice = 0 OR p.price >= targetPrice)
    			AND (targetRating = 0 OR p.rating >= targetRating)
			)
			LIMIT targetLimit 
    		OFFSET targetOffset;
			End If;
END;
$$ LANGUAGE plpgsql;

select * from getProducts_fn(null,'xiomi','',0,0,10,0)

