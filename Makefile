f1=real_entry
f2=entry
f3=cipher
bin=cypher
test:
	echo "something" > $(f1)
	cp $(f1) $(f2)
	diff $(f1) $(f2)
	./$(bin) e
	rm $(f2)
	./$(bin) d
	rm $(f3)
	diff $(f1) $(f2)
	rm $(f1) $(f2)
