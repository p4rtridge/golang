CREATE OR REPLACE FUNCTION update_order_total_price()
RETURNS TRIGGER AS $$

DECLARE
  new_total_price float;

BEGIN
  SELECT COALESCE(SUM(product_price * quantity), 0) INTO new_total_price FROM order_items WHERE order_id = NEW.order_id;
  UPDATE orders SET total_price = new_total_price WHERE id = NEW.order_id;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trigger_update_order_total_price
AFTER INSERT OR UPDATE ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_order_total_price();
