insert into tags (id, name) values 
(1, 'muebles'),
(2, 'oficina'),
(3, 'cocina'),
(4, 'herramientas');

insert into articles_conditions (id, slug, label, description) values
(1, 'new', 'Nuevo', 'el articulo se compro, pero nunca se uso si no fue para probarlo'),
(2, 'semi-new', 'Semi nuevo', 'el articulo se uso solo un par de veces y esta en perfecto estado'),
(3, 'used', 'Usado', 'el articulo tuvo un uso normal y muestra signos de uso'),
(4, 'very-used', 'Muy usado', 'el articulo se uso mucho, o por mucho tiempo, y muestra signos de deterioro');

insert into articles (id, name, description, reference_price, condition_id, available_for_trade) values
(1, 'escritorio 50x120 melamina', 'escritorio comprado en sodimac', 49990, 1, false),
(2, 'olla roja', 'se ve fea pero esta en buen estado. La herede y no la uso', 0, 4, true),
(3, 'caladora black&decker alambrica', 'perfecto estado, con sus accesorios y caja. Se vende por upgrade', 25990, 3, true),
(4, 'Pantalla plana lcd viewsonic 4:3', 'viejita y pesada pero sirve para tenerla de pantalla auxiliar par un servidor, escritorio secundario, o algo por el estilo', 0, 4, true);

insert into articles_tags (article_id, tag_id) values
(1, 1),
(1, 2),
(2, 3),
(3, 4),
(4, 2);

insert into wishitems (id, name, observed_price, external_url) values
(1, 'Cámara Web Pc Notebook Full Hd Usb Microfono 1080p', 39990, 'https://articulo.mercadolibre.cl/MLC-2973994090-camara-web-pc-notebook-full-hd-usb-microfono-1080p-_JM#polycard_client=recommendations_vip-pads-right&reco_backend=ranker_retrieval_system_ads&reco_model=rk_ent_v3_retsys_ads&reco_client=vip-pads-right&reco_item_pos=1&reco_backend_type=low_level&reco_id=37cce21a-fec0-4350-b3c1-781491d263d7&is_advertising=true&ad_domain=VIPCORE_RIGHT&ad_position=2&ad_click_id=M2M2ZDYyMDQtYzBmNC00ZTY2LTk3N2EtOTY5NDFmOTM5ZDAw'),
(2, 'Advantage360 Signature Series', 599000, 'https://kinesis-ergo.com/shop/advantage360-signature/');

insert into wishitems_tags (wishitem_id, tag_id) values
(1, 2),
(2, 1);
