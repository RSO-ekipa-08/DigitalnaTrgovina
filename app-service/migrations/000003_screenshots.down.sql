UPDATE applications a
SET screenshots = CASE a.id
    WHEN 'ec9d81fd-551b-4fdc-b220-3712b241a49a' THEN ARRAY['ft_screen1.jpg', 'ft_screen2.jpg', 'ft_screen3.jpg']
    WHEN '9d73b8fc-7876-4d26-89a0-ca10d1a93457' THEN ARRAY['puzzle1.jpg', 'puzzle2.jpg']
    WHEN '4d6a2283-4b76-4359-bba1-f7dc88b6f2b1' THEN ARRAY['task1.jpg', 'task2.jpg', 'task3.jpg']
    WHEN '8960306c-3588-450f-abdc-11bedf81dfe8' THEN ARRAY['lang1.jpg', 'lang2.jpg', 'lang3.jpg']
    WHEN 'c23a3d46-8ffd-4183-a5f8-10bd8d54b6c3' THEN ARRAY['social1.jpg', 'social2.jpg']
    WHEN '1e68d972-5318-402c-a3e8-c8e45d9e5eb9' THEN ARRAY['movie1.jpg', 'movie2.jpg']
    WHEN '8d79c9a6-512a-460a-b43c-90244b615d5d' THEN ARRAY['diet1.jpg', 'diet2.jpg']
    WHEN '665d41a1-cd27-4473-a56a-80de5f6ef6e4' THEN ARRAY['invoice1.jpg', 'invoice2.jpg']
    WHEN 'f78a6f0c-b0d3-4ec4-91f6-df5e3162704e' THEN ARRAY['files1.jpg', 'files2.jpg']
    WHEN 'f951d41d-dec1-40df-8c45-ce517f054d1e' THEN ARRAY['travel1.jpg', 'travel2.jpg']
    WHEN 'bceee33e-32e2-4f0d-930b-90e5daf8be24' THEN ARRAY['shop1.jpg', 'shop2.jpg']
    WHEN 'bab4443a-49f0-4cd9-a90f-fd72e0e7f1ef' THEN ARRAY['math1.jpg', 'math2.jpg']
    WHEN '1ddaebd7-548d-4976-a201-e2a8154ddfd4' THEN ARRAY['race1.jpg', 'race2.jpg', 'race3.jpg']
    WHEN '5fb9fc55-1ec6-42f9-a75c-291419122e74' THEN ARRAY['med1.jpg', 'med2.jpg']
    WHEN 'c206536d-a8ac-4e10-8f0e-4eafc063dab9' THEN ARRAY['photo1.jpg', 'photo2.jpg']
    WHEN 'd6c1734f-f451-4833-99c7-1a8dda2c4305' THEN ARRAY['weather1.jpg', 'weather2.jpg']
    WHEN '0400e28f-d314-4b0c-a136-5714a7854125' THEN ARRAY['recipe1.jpg', 'recipe2.jpg']
    WHEN '62efb487-922c-40d6-bec1-2ab265b9c488' THEN ARRAY['budget1.jpg', 'budget2.jpg']
    WHEN 'a140970c-8ca7-4102-be63-324c2ce0a49d' THEN ARRAY['music1.jpg', 'music2.jpg']
    WHEN '2f0c4159-55f0-4009-a309-7011423de22d' THEN ARRAY['cal1.jpg', 'cal2.jpg']
    ELSE screenshots
END;
